package ai

import (
	"strconv"
	"strings"
	"time"

	"github.com/tsybulin/gosha/aut"
	"github.com/tsybulin/gosha/evt"
)

type greaterThanTrigger struct {
	aut.Trigger
	entityID     string
	value        float64
	firstFiredAt time.Time
	forSeconds   time.Duration
}

func (t *greaterThanTrigger) GetEntityID() string {
	return t.entityID
}

func (t *greaterThanTrigger) GetValue() float64 {
	return t.value
}

func (t *greaterThanTrigger) FireCompare(event evt.Event) bool {
	if event.EventType != "state_changed" {
		return false
	}

	if t.GetEntityID() != event.Data.EntityID {
		return false
	}

	if value, err := strconv.ParseFloat(event.Data.NewState.State, 64); err == nil {
		rslt := value > t.value

		if t.forSeconds <= 0 {
			return rslt
		}

		if !rslt {
			t.firstFiredAt = time.Time{}
			return false
		}

		now := time.Now()

		if t.firstFiredAt.IsZero() {
			t.firstFiredAt = now
			return false
		}

		if now.Sub(t.firstFiredAt) >= t.forSeconds*time.Second {
			t.firstFiredAt = time.Time{}
			return true
		}
	}

	return false
}

func newGreaterThanTrigger(id string, value float64, forSeconds time.Duration) aut.CompareTrigger {
	return &greaterThanTrigger{
		Trigger:      newTrigger("greater_than"),
		entityID:     id,
		value:        value,
		firstFiredAt: time.Time{},
		forSeconds:   forSeconds,
	}
}

func newGreaterThanTriggers(cfg map[string]string) []aut.CompareTrigger {
	triggers := make([]aut.CompareTrigger, 0)

	for _, id := range strings.Fields(cfg["components"]) {
		if value, err := strconv.ParseFloat(cfg["value"], 64); err == nil {
			forSeconds, _ := strconv.ParseInt(cfg["for"], 10, 64)
			triggers = append(triggers, newGreaterThanTrigger(id, value, time.Duration(forSeconds)))
		}
	}

	return triggers
}
