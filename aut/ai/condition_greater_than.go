package ai

import (
	"strconv"
	"strings"
	"time"

	"github.com/tsybulin/gosha/aut"
)

type greaterThanCondition struct {
	aut.Condition
	entityID         string
	value            float64
	firstSatisfiedAt time.Time
	forSeconds       time.Duration
}

func (gtc *greaterThanCondition) GetEntityID() string {
	return gtc.entityID
}

func (gtc *greaterThanCondition) GetValue() float64 {
	return gtc.value
}

func (gtc *greaterThanCondition) SatisfiedCompare(id string, value float64) bool {
	if gtc.GetEntityID() != id {
		return false
	}

	rslt := value > gtc.value

	if gtc.forSeconds <= 0 {
		return rslt
	}

	if !rslt {
		gtc.firstSatisfiedAt = time.Time{}
		return false
	}

	now := time.Now()

	if gtc.firstSatisfiedAt.IsZero() {
		gtc.firstSatisfiedAt = now
		return false
	}

	if now.Sub(gtc.firstSatisfiedAt) >= gtc.forSeconds*time.Second {
		gtc.firstSatisfiedAt = time.Time{}
		return true
	}

	return false
}

func newgGeaterThanCondition(id string, value float64, forSeconds time.Duration) aut.CompareCondition {
	return &greaterThanCondition{
		Condition:        newCondition("greater_than"),
		entityID:         id,
		value:            value,
		firstSatisfiedAt: time.Time{},
		forSeconds:       forSeconds,
	}
}

func newGreaterThanConditions(cfg map[string]string) []aut.CompareCondition {
	conditions := make([]aut.CompareCondition, 0)

	for _, id := range strings.Fields(cfg["components"]) {
		if value, err := strconv.ParseFloat(cfg["value"], 64); err == nil {
			forSeconds, _ := strconv.ParseInt(cfg["for"], 10, 64)
			conditions = append(conditions, newgGeaterThanCondition(id, value, time.Duration(forSeconds)))
		}
	}
	return conditions
}
