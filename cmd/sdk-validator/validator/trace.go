package validator

import (
	"context"
	"fmt"

	tracev1 "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"go.uber.org/multierr"
	"google.golang.org/protobuf/encoding/protojson"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/otelcollector"
)

// TraceHandler implements ExportTraceService.
type TraceHandler struct {
	tracev1.UnimplementedTraceServiceServer
}

// Export is a dummy Export handler.
func (t TraceHandler) Export(ctx context.Context, req *tracev1.ExportTraceServiceRequest) (*tracev1.ExportTraceServiceResponse, error) {
	log.Trace().Msg("Received Export request")

	var merr error

	for _, resourceSpans := range req.ResourceSpans {
		for _, scopeSpan := range resourceSpans.ScopeSpans {
			for _, span := range scopeSpan.Spans {
				var err error
				var flowStartTS, flowEndTS, workloadTS int64
				for _, attribute := range span.Attributes {
					switch attribute.Key {
					case otelcollector.ApertureCheckResponseLabel:
						log.Trace().Str("attribute", otelcollector.ApertureCheckResponseLabel).Msg("Validating attribute")
						v := attribute.Value.GetStringValue()
						checkResponse := &flowcontrolv1.CheckResponse{}
						perr := protojson.Unmarshal([]byte(v), checkResponse)
						if perr != nil {
							log.Error().Err(perr).Msg("Failed to validate flowcontrol CheckResponse")
							err = multierr.Append(err, fmt.Errorf("invalid %s: %w", otelcollector.ApertureCheckResponseLabel, perr))
						}
					case otelcollector.ApertureSourceLabel:
						log.Trace().Str("attribute", otelcollector.ApertureSourceLabel).Msg("Validating attribute")
						v := attribute.Value.GetStringValue()
						if v != "sdk" {
							log.Error().Msg("Failed to validate source")
							err = multierr.Append(err, fmt.Errorf("invalid %s", otelcollector.ApertureSourceLabel))
						}
					case otelcollector.ApertureFeatureStatusLabel:
						log.Trace().Str("attribute", otelcollector.ApertureFeatureStatusLabel).Msg("Validating attribute")
						v := attribute.Value.GetStringValue()
						if v != otelcollector.ApertureFeatureStatusOK && v != otelcollector.ApertureFeatureStatusError {
							log.Error().Msg("Failed to validate feature status")
							err = multierr.Append(err, fmt.Errorf("invalid %s", otelcollector.ApertureFeatureStatusLabel))
						}
					case otelcollector.ApertureFlowStartTimestampLabel:
						flowStartTS = attribute.Value.GetIntValue()
					case otelcollector.ApertureFlowEndTimestampLabel:
						flowEndTS = attribute.Value.GetIntValue()
					case otelcollector.ApertureWorkloadStartTimestampLabel:
						workloadTS = attribute.Value.GetIntValue()
					}
				}
				log.Trace().Str("attribute", otelcollector.ApertureFlowStartTimestampLabel).Str("attribute", otelcollector.ApertureFlowEndTimestampLabel).Msg("Validating attribute")
				if flowStartTS > flowEndTS {
					log.Error().Msg("Failed to validate start and end flow timestamps")
					err = multierr.Append(err, fmt.Errorf("invalid %s and %s", otelcollector.ApertureFlowStartTimestampLabel, otelcollector.ApertureFlowEndTimestampLabel))
				}
				log.Trace().Str("attribute", otelcollector.ApertureWorkloadStartTimestampLabel).Msg("Validating attribute")
				if workloadTS == 0 {
					log.Error().Msg("Failed to validate workload start timestamp")
					err = multierr.Append(err, fmt.Errorf("invalid %s", otelcollector.ApertureWorkloadStartTimestampLabel))
				}
				merr = multierr.Append(merr, err)
				if merr != nil {
					return &tracev1.ExportTraceServiceResponse{}, merr
				}
			}
		}
	}

	log.Info().Msg("Validated span attributes")
	return &tracev1.ExportTraceServiceResponse{}, merr
}
