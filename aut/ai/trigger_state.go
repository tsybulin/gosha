package ai

import (
	"strings"

	"github.com/tsybulin/gosha/aut"
	"github.com/tsybulin/gosha/evt"
)

type stateTrigger struct {
	aut.Trigger
	entityID string
	from     string
	to       string
}

func (t *stateTrigger) GetEntityID() string {
	return t.entityID
}

func (t *stateTrigger) GetFrom() string {
	return t.from
}

func (t *stateTrigger) GetTo() string {
	return t.to
}

func (t *stateTrigger) FireState(event evt.Event) bool {
	if event.EventType != "state_changed" {
		return false
	}

	if t.GetEntityID() != event.Data.EntityID {
		return false
	}

	if len(t.GetFrom()) > 0 {
		if event.Data.OldState == nil || event.Data.OldState.State != t.GetFrom() {
			return false
		}
	}

	if len(t.GetTo()) > 0 {
		if event.Data.NewState == nil || event.Data.NewState.State != t.GetTo() {
			return false
		}
	}

	return true
}

func newStateTrigger(id string, from, to string) aut.StateTrigger {
	return &stateTrigger{
		Trigger:  newTrigger("state"),
		entityID: id,
		from:     from,
		to:       to,
	}
}

func newStateTriggers(cfg map[string]string) []aut.StateTrigger {
	triggers := make([]aut.StateTrigger, 0)

	for _, id := range strings.Fields(cfg["components"]) {
		triggers = append(triggers, newStateTrigger(id, cfg["from"], cfg["to"]))
	}

	return triggers
}
