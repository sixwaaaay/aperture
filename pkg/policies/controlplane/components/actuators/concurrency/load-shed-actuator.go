package concurrency

import (
	"context"
	"path"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"

	policydecisionsv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/decisions/v1"
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	wrappersv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/wrappers/v1"
	"github.com/fluxninja/aperture/pkg/config"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdwriter "github.com/fluxninja/aperture/pkg/etcd/writer"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/common"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
	"github.com/rs/zerolog"
)

// LoadShedActuator struct.
type LoadShedActuator struct {
	policyReadAPI  iface.Policy
	decision       *policydecisionsv1.LoadShedDecision
	etcdPath       string
	writer         *etcdwriter.Writer
	agentGroupName string
	componentIndex int
}

// NewLoadShedActuatorAndOptions creates load shed actuator and its fx options.
func NewLoadShedActuatorAndOptions(
	_ *policylangv1.LoadShedActuator,
	componentIndex int,
	policyReadAPI iface.Policy,
	agentGroupName string,
) (runtime.Component, fx.Option, error) {
	etcdPath := path.Join(common.LoadShedDecisionsPath,
		common.DataplaneComponentKey(agentGroupName, policyReadAPI.GetPolicyName(), int64(componentIndex)))
	lsa := &LoadShedActuator{
		policyReadAPI:  policyReadAPI,
		agentGroupName: agentGroupName,
		componentIndex: componentIndex,
		etcdPath:       etcdPath,
	}
	lsa.decision = &policydecisionsv1.LoadShedDecision{}

	return lsa, fx.Options(
		fx.Invoke(lsa.setupWriter),
	), nil
}

func (lsa *LoadShedActuator) setupWriter(etcdClient *etcdclient.Client, lifecycle fx.Lifecycle) error {
	logger := lsa.policyReadAPI.GetStatusRegistry().GetLogger()
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			lsa.writer = etcdwriter.NewWriter(etcdClient, true)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			_, err := etcdClient.KV.Delete(clientv3.WithRequireLeader(ctx), lsa.etcdPath)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to delete load shed decision config")
				return err
			}
			lsa.writer.Close()
			return nil
		},
	})

	return nil
}

// Execute implements runtime.Component.Execute.
func (lsa *LoadShedActuator) Execute(inPortReadings runtime.PortToValue, tickInfo runtime.TickInfo) (runtime.PortToValue, error) {
	// Get the decision from the port
	lsf, ok := inPortReadings["load_shed_factor"]
	if ok {
		if len(lsf) > 0 {
			lsfReading := lsf[0]
			var lsfValue float64
			if !lsfReading.Valid() {
				lsfValue = 0
			} else {
				if lsfReading.Value() <= 0 {
					lsfValue = 0
				} else {
					lsfValue = lsfReading.Value()
				}
			}
			return nil, lsa.publishLoadShedFactor(lsfValue)
		}
	}
	return nil, nil
}

// DynamicConfigUpdate is a no-op for load shed actuator.
func (lsa *LoadShedActuator) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
}

func (lsa *LoadShedActuator) publishLoadShedFactor(loadShedFactor float64) error {
	logger := lsa.policyReadAPI.GetStatusRegistry().GetLogger()
	// Publish only if there's a change
	if lsa.decision.GetLoadShedFactor() != loadShedFactor {
		// Save load shed factor in decision message
		lsa.decision.LoadShedFactor = loadShedFactor
		// Publish decision
		logger.Sample(zerolog.Often).Debug().Float64("loadShedFactor", loadShedFactor).Msg("Publish load shed decision")
		wrapper := &wrappersv1.LoadShedDecisionWrapper{
			LoadShedDecision: lsa.decision,
			CommonAttributes: &wrappersv1.CommonAttributes{
				PolicyName:     lsa.policyReadAPI.GetPolicyName(),
				PolicyHash:     lsa.policyReadAPI.GetPolicyHash(),
				ComponentIndex: int64(lsa.componentIndex),
			},
		}
		dat, err := proto.Marshal(wrapper)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to marshal policy decision")
			return err
		}
		lsa.writer.Write(lsa.etcdPath, dat)
	}
	return nil
}
