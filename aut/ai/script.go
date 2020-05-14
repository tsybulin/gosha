package ai

import (
	"github.com/tsybulin/gosha/aut"
	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/evt"
)

type script struct {
	id      string
	state   string
	actions []aut.Action
}

func (s *script) GetID() string {
	return s.id
}

func (s *script) IsActive() bool {
	return s.state == "active"
}

func (s *script) SetActive(active bool) {
	if active {
		s.state = "active"
	} else {
		s.state = "inactive"
	}
}

func (s *script) GetState() evt.State {
	event := evt.State{
		EntityID:   s.GetID(),
		State:      s.state,
		Attributes: make(map[string]interface{}, 0),
	}

	event.Attributes["actions"] = len(s.actions)

	return event
}

func (s *script) Execute() {
	for _, a := range s.actions {
		a.Execute()
	}
}

// NewScripts ...
func NewScripts(scfg []struct {
	Script  string
	Actions []map[string]string
}) []aut.Script {
	scripts := make([]aut.Script, 0)

	for _, sc := range scfg {
		s := &script{
			id:      cmp.DomainScript.String() + "." + sc.Script,
			state:   "inactive",
			actions: make([]aut.Action, 0),
		}

		for _, acc := range sc.Actions {
			for _, ac := range newActions(acc) {
				s.actions = append(s.actions, ac)
			}
		}

		scripts = append(scripts, s)
	}

	return scripts
}
