package iface

import (
	"context"
	"strconv"

	selectorv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/selector/v1"
	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	"github.com/prometheus/client_golang/prometheus"
)

//go:generate mockgen -source=limiter.go -destination=../../mocks/mock_limiter.go -package=mocks

// LimiterID is the ID of the Limiter.
type LimiterID struct {
	PolicyName     string
	PolicyHash     string
	ComponentIndex int64
}

// String function returns the LimiterID as a string.
func (limiterID LimiterID) String() string {
	return "policy_name-" + limiterID.PolicyName + "-component_index-" + strconv.FormatInt(limiterID.ComponentIndex, 10) + "-policy_hash-" + limiterID.PolicyHash
}

// Limiter interface.
// Lifetime of this interface is per policy/component.
type Limiter interface {
	GetPolicyName() string
	GetSelector() *selectorv1.Selector
	RunLimiter(ctx context.Context, labels map[string]string) *flowcontrolv1.LimiterDecision
	GetLimiterID() LimiterID
}

// RateLimiter interface.
type RateLimiter interface {
	Limiter
	TakeN(labels map[string]string, count int) (label string, ok bool, remaining int, current int)
	GetCounter() prometheus.Counter
}

// ConcurrencyLimiter interface.
type ConcurrencyLimiter interface {
	Limiter
	GetObserver(labels map[string]string) prometheus.Observer
}
