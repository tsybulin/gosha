package svc

import (
	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/evt"
	"github.com/tsybulin/gosha/logger"
)

// Group ...
type Group interface {
	Service
	TurnOn(string)
	TurnOff(string)
	Toggle(string)
}

type groupService struct {
	Service
	eventBus evt.Bus
	groups   map[string]cmp.Group
}

func (gs *groupService) Components() []string {
	components := make([]string, 0)
	for id := range gs.groups {
		components = append(components, id)
	}
	return components
}

func (gs *groupService) Description() Description {
	return Description{
		ID: gs.GetID(),
		Methods: []Method{
			{Method: "turn_on", Parameters: []string{"group_id"}},
			{Method: "turn_off", Parameters: []string{"group_id"}},
			{Method: "toggle", Parameters: []string{"group_id"}},
		},
	}
}

func (gs *groupService) States() StateResult {
	stateResults := make([]evt.State, 0)

	for _, c := range gs.groups {
		stateResults = append(stateResults, c.GetState())
	}

	return StateResult{
		ID:      1,
		Type:    "result",
		Success: true,
		Result:  stateResults,
	}
}

func (gs *groupService) registerGroup(g cmp.Group) {
	if g.GetDomain() == cmp.DomainGroup {
		gs.groups[g.GetID()] = g

		go func() {
			gs.eventBus.Publish(logger.Topic, logger.LevelDebug, "GroupService.registerGroup", g.GetID())
		}()
	}
}

func (gs *groupService) registerComponent(c cmp.Component) {
	for _, g := range gs.groups {
		if g.AddComponent(c) {
			go func() {
				gs.eventBus.Publish(logger.Topic, logger.LevelDebug, "GroupService.registerComponent group:", g.GetID(), "component:", c.GetID())
			}()
		}
	}
}

func (gs *groupService) TurnOn(id string) {
	if g := gs.groups[id]; g != nil {
		for _, c := range g.GetComponents() {
			if s, ok := c.(cmp.Switch); ok && !s.IsOn() {
				gs.eventBus.Publish(evt.StateChangedTopic, eventFor(s, s.TurnOn))
				gs.eventBus.Publish(evt.PtfPushTopic, c, "power")
			}
		}
	} else {
		gs.eventBus.Publish(logger.Topic, logger.LevelWarn, "GroupService.TurnOn unknown id:", id)
	}
}

func (gs *groupService) TurnOff(id string) {
	if g := gs.groups[id]; g != nil {
		for _, c := range g.GetComponents() {
			if s, ok := c.(cmp.Switch); ok && s.IsOn() {
				gs.eventBus.Publish(evt.StateChangedTopic, eventFor(s, s.TurnOff))
				gs.eventBus.Publish(evt.PtfPushTopic, c, "power")
			}
		}
	} else {
		gs.eventBus.Publish(logger.Topic, logger.LevelWarn, "GroupService.TurnOff unknown id:", id)
	}
}

func (gs *groupService) Toggle(id string) {
	if g := gs.groups[id]; g != nil {
		for _, c := range g.GetComponents() {
			if s, ok := c.(cmp.Switch); ok {
				gs.eventBus.Publish(evt.StateChangedTopic, eventFor(s, s.Toggle))
				gs.eventBus.Publish(evt.PtfPushTopic, c, "power")
			}
		}
	} else {
		gs.eventBus.Publish(logger.Topic, logger.LevelWarn, "GroupService.Toggle unknown id:", id)
	}
}

func newGroupService(eventBus evt.Bus) Group {
	gs := &groupService{
		Service:  newService(cmp.DomainGroup),
		eventBus: eventBus,
		groups:   make(map[string]cmp.Group, 0),
	}

	gs.eventBus.Subscribe(cmp.RegisterGroupTopic, "GroupService.registerGroup", gs.registerGroup)
	gs.eventBus.Subscribe(cmp.RegisterTopic, "GroupService.registerComponent", gs.registerComponent)

	return gs
}
