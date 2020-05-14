package ai

import (
	"strings"
	"time"

	"github.com/tsybulin/gosha/aut"
)

type timeTrigger struct {
	aut.Trigger
	at    time.Time
	fired bool
}

func (tt *timeTrigger) GetAt() time.Time {
	return tt.at
}

func (tt *timeTrigger) FireTime(t time.Time) bool {
	eq := tt.at.Hour() == t.Hour() &&
		tt.at.Minute() == t.Minute()

	if !eq {
		tt.fired = false
		return tt.fired
	}

	if tt.fired {
		return false
	}

	tt.fired = true
	return true
}

func newTimeTriggerAt(at time.Time) aut.TimeTrigger {
	return &timeTrigger{
		Trigger: newTrigger("time"),
		at:      at,
		fired:   false,
	}
}

func newTimeTriggers(cfg map[string]string) []aut.TimeTrigger {
	triggers := make([]aut.TimeTrigger, 0)

	for _, at := range strings.Fields(cfg["at"]) {
		tm, err := time.Parse("15:04", at)
		if err != nil {
			tm = time.Now()
		}
		triggers = append(triggers, newTimeTriggerAt(tm))
	}

	return triggers
}
