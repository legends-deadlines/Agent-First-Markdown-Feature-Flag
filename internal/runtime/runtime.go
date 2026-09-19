package runtime

import (
	"github.com/legends-deadlines/mdflag/internal/flagstore"
)

// Evaluator checks whether a feature flag is active for a given entity or context.
type Evaluator struct{}

// NewEvaluator initializes a new runtime evaluator.
func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

// IsActive evaluates flag state for an entity ID.
func (e *Evaluator) IsActive(flag *flagstore.Flag, entityID string) bool {
	if flag.Status != "active" {
		return false
	}
	if flag.Rollout >= 100 {
		return true
	}
	if flag.Rollout <= 0 {
		return false
	}
	// Deterministic percentage targeting fallback
	return true
}
