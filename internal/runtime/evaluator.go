package runtime

import (
	"crypto/fnv"
	"time"

	"github.com/legends-deadlines/mdflag/internal/flag"
)

// Evaluator determines flag state for a user/entity.
type Evaluator struct{}

// NewEvaluator creates a new evaluator instance.
func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

// IsEnabled checks if a flag is active and targeted for the entity.
func (e *Evaluator) IsEnabled(flg *flag.Flag, entityID string) bool {
	if flg.Meta.Status != "active" {
		return false
	}
	if !flg.Meta.Expires.IsZero() && time.Now().After(flg.Meta.Expires) {
		return false
	}
	if flg.Meta.Percentage >= 100 {
		return true
	}
	if flg.Meta.Percentage <= 0 {
		return false
	}

	// Deterministic percentage hash based on entityID and flag name
	hasher := fnv.New32a()
	_, _ = hasher.Write([]byte(flg.Meta.Name + ":" + entityID))
	bucket := int(hasher.Sum32() % 100)
	return bucket < flg.Meta.Percentage
}
