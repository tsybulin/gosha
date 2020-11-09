package ai

import (
	"strconv"
	"strings"
	"time"

	"github.com/tsybulin/gosha/aut"
)

type lessThanCondition struct {
	aut.Condition
	entityID         string
	value            float64
	firstSatisfiedAt time.Time
	forSeconds       time.Duration
}

func (ltc *lessThanCondition) GetEntityID() string {
	return ltc.entityID
}

func (ltc *lessThanCondition) GetValue() float64 {
	return ltc.value
}

func (ltc *lessThanCondition) SatisfiedCompare(id string, value float64) bool {
	if ltc.GetEntityID() != id {
		return false
	}

	rslt := value < ltc.value

	if ltc.forSeconds <= 0 {
		return rslt
	}

	if !rslt {
		ltc.firstSatisfiedAt = time.Time{}
		return false
	}

	now := time.Now()

	if ltc.firstSatisfiedAt.IsZero() {
		ltc.firstSatisfiedAt = now
		return false
	}

	if now.Sub(ltc.firstSatisfiedAt) >= ltc.forSeconds*time.Second {
		ltc.firstSatisfiedAt = time.Time{}
		return true
	}

	return false
}

func newLessThanCondition(id string, value float64, forSeconds time.Duration) aut.CompareCondition {
	return &lessThanCondition{
		Condition:        newCondition("less_than"),
		entityID:         id,
		value:            value,
		firstSatisfiedAt: time.Time{},
		forSeconds:       forSeconds,
	}
}

func newLessThanConditions(cfg map[string]string) []aut.CompareCondition {
	conditions := make([]aut.CompareCondition, 0)

	for _, id := range strings.Fields(cfg["components"]) {
		if value, err := strconv.ParseFloat(cfg["value"], 64); err == nil {
			forSeconds, _ := strconv.ParseInt(cfg["for"], 10, 64)
			conditions = append(conditions, newLessThanCondition(id, value, time.Duration(forSeconds)))
		}
	}

	return conditions
}
