package svc

import (
	"github.com/tsybulin/gosha/aut"
	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/evt"
	"github.com/tsybulin/gosha/logger"
)

// Script ...
type Script interface {
	Service
	Execute(string)
}

type scriptService struct {
	Service
	eventBus evt.Bus
	scripts  map[string]aut.Script
}

func (ss *scriptService) Components() []string {
	components := make([]string, 0)
	for id := range ss.scripts {
		components = append(components, id)
	}
	return components
}

func (ss *scriptService) Description() Description {
	return Description{
		ID: ss.GetID(),
		Methods: []Method{
			{Method: "execute", Parameters: []string{"script_id"}},
		},
	}
}

func (ss *scriptService) States() StateResult {
	stateResults := make([]evt.State, 0)

	for _, c := range ss.scripts {
		stateResults = append(stateResults, c.GetState())
	}

	return StateResult{
		ID:      1,
		Type:    "result",
		Success: true,
		Result:  stateResults,
	}
}

func (ss *scriptService) Execute(id string) {
	go func() {
		if s := ss.scripts[id]; s != nil {
			ss.eventBus.Publish(logger.Topic, logger.LevelInfo, "ScriptService.Execute", id)

			event := eventFor(s, func() {
				s.SetActive(true)
			})

			ss.eventBus.Publish(evt.StateChangedTopic, event)

			s.Execute()

			event = eventFor(s, func() {
				s.SetActive(false)
			})

			ss.eventBus.Publish(evt.StateChangedTopic, event)
		} else {
			ss.eventBus.Publish(logger.Topic, logger.LevelWarn, "ScriptService.Execute unknown id:", id)
		}
	}()
}

func (ss *scriptService) registerScript(s aut.Script) {
	ss.scripts[s.GetID()] = s
	go func() {
		ss.eventBus.Publish(logger.Topic, logger.LevelDebug, "ScriptService.registerScript", s.GetID())
	}()
}

func newScriptService(eventBus evt.Bus) Script {
	ss := &scriptService{
		Service:  newService(cmp.DomainScript),
		eventBus: eventBus,
		scripts:  make(map[string]aut.Script, 0),
	}

	ss.eventBus.Subscribe(aut.ScriptRegisterTopic, "ScriptService.registerScript", ss.registerScript)

	return ss
}
