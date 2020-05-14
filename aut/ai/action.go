package ai

import (
	"strings"

	"github.com/tsybulin/gosha/aut"
	"github.com/tsybulin/gosha/svc"
)

type action struct {
	service   string
	action    string
	component string
}

func (a *action) GetService() string {
	return a.service
}

func (a *action) GetAction() string {
	return a.action
}

func (a *action) GetComponent() string {
	return a.component
}

func (a *action) Execute() {
	go func() {
		svc.NewRegistry(nil).Execute(a.GetService(), a.GetAction(), a.GetComponent())
	}()
}

func newAction(service, act, component string) aut.Action {
	return &action{
		service:   service,
		action:    act,
		component: component,
	}
}

func newActions(cfg map[string]string) []aut.Action {
	actions := make([]aut.Action, 0)
	for _, component := range strings.Fields(cfg["components"]) {
		actions = append(actions, newAction(cfg["service"], cfg["action"], component))
	}
	return actions
}
