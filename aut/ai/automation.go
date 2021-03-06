package ai

import (
	"sync"

	"github.com/tsybulin/gosha/aut"
)

type automation struct {
	id         string
	triggers   []aut.Trigger
	conditions []aut.Condition
	actions    []aut.Action
	mux        sync.Mutex
}

func (a *automation) GetID() string {
	return a.id
}

func (a *automation) GetTriggers() []aut.Trigger {
	return a.triggers
}

func (a *automation) GetActions() []aut.Action {
	return a.actions
}

func (a *automation) GetContitions() []aut.Condition {
	return a.conditions
}

func (a *automation) Lock() {
	a.mux.Lock()
}

func (a *automation) Unlock() {
	a.mux.Unlock()
}

// NewAutomations ...
func NewAutomations(acfg []struct {
	Automation string
	Triggers   []map[string]string
	Conditions []map[string]string
	Actions    []map[string]string
}) []aut.Automation {

	aau := make([]aut.Automation, 0)

	for _, auc := range acfg {
		au := &automation{
			id:         auc.Automation,
			triggers:   make([]aut.Trigger, 0),
			conditions: make([]aut.Condition, 0),
			actions:    make([]aut.Action, 0),
		}

		for _, trc := range auc.Triggers {
			switch trc["platform"] {
			case "state":
				for _, tr := range newStateTriggers(trc) {
					au.triggers = append(au.triggers, tr)
				}
			case "time":
				for _, tr := range newTimeTriggers(trc) {
					au.triggers = append(au.triggers, tr)
				}
			case "less_than":
				for _, tr := range newLessThanTriggers(trc) {
					au.triggers = append(au.triggers, tr)
				}
			case "greater_than":
				for _, tr := range newGreaterThanTriggers(trc) {
					au.triggers = append(au.triggers, tr)
				}
			}
		}

		for _, coc := range auc.Conditions {
			switch coc["platform"] {
			case "state":
				for _, co := range newStateConditions(coc) {
					au.conditions = append(au.conditions, co)
				}
			case "less_than":
				for _, co := range newLessThanConditions(coc) {
					au.conditions = append(au.conditions, co)
				}
			case "greater_than":
				for _, co := range newGreaterThanConditions(coc) {
					au.conditions = append(au.conditions, co)
				}
			}
		}

		for _, acc := range auc.Actions {
			for _, ac := range newActions(acc) {
				au.actions = append(au.actions, ac)
			}
		}

		aau = append(aau, au)
	}

	return aau
}
