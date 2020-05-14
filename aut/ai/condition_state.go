package ai

import (
	"strings"

	"github.com/tsybulin/gosha/aut"
)

type stateCondition struct {
	aut.Condition
	entityID string
	state    string
}

func (c *stateCondition) GetEntityID() string {
	return c.entityID
}

func (c *stateCondition) GetState() string {
	return c.state
}

func (c *stateCondition) SatisfiedState(id, state string) bool {
	if c.GetEntityID() != id {
		return false
	}

	return c.GetState() == state
}

func newStateCondition(id, state string) aut.StateCondition {
	return &stateCondition{
		Condition: newCondition("state"),
		entityID:  id,
		state:     state,
	}
}

func newStateConditions(cfg map[string]string) []aut.StateCondition {
	conditions := make([]aut.StateCondition, 0)

	for _, id := range strings.Fields(cfg["components"]) {
		conditions = append(conditions, newStateCondition(id, cfg["state"]))
	}

	return conditions
}
