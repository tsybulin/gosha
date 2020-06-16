package aut

import (
	"fmt"
	"sync"
	"time"

	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/evt"
	"github.com/tsybulin/gosha/logger"
)

// Automator ...
type Automator interface {
}

type automator struct {
	eventBus    evt.Bus
	automations map[string]Automation
	exectimes   map[string]time.Time
	components  map[string]cmp.Component
}

func (ar *automator) registerAutomation(a Automation) {
	ar.automations[a.GetID()] = a
	ar.exectimes[a.GetID()] = time.Now()
	go func() {
		ar.eventBus.Publish(logger.Topic, logger.LevelInfo, "Automator.registerAutomation", a.GetID())
	}()
}

func (ar *automator) registerComponent(c cmp.Component) {
	ar.components[c.GetID()] = c
}

func (ar *automator) executeIfConditions(au Automation) {
	fired := true

	for _, co := range au.GetContitions() {
		switch co.(type) {
		case StateCondition:
			sc := co.(StateCondition)
			if c := ar.components[sc.GetEntityID()]; c != nil {
				switch c.(type) {
				case cmp.Switchable:
					fired = fired && sc.SatisfiedState(c.GetID(), c.(cmp.Switchable).GetOnString())
				case cmp.BinarySensor:
					fired = fired && sc.SatisfiedState(c.GetID(), c.(cmp.BinarySensor).GetOnString())
				case cmp.Sensor:
					fired = fired && sc.SatisfiedState(c.GetID(), fmt.Sprint(c.(cmp.Sensor).GetValue()))
				case cmp.Timer:
					fired = fired && sc.SatisfiedState(c.GetID(), c.(cmp.Timer).GetTimerState().String())
				}
			}

		}

		if !fired {
			break
		}
	}

	if !fired {
		return
	}

	ar.eventBus.Publish(logger.Topic, logger.LevelSystem, "Automator.Execute ", au.GetID())
	ar.exectimes[au.GetID()] = time.Now()

	for _, ac := range au.GetActions() {
		ac.Execute()
		ar.eventBus.Publish(logger.Topic, logger.LevelSystem, "Automator.Execute action", au.GetID(), ac.GetService(), ac.GetAction(), ac.GetComponent())
	}
}

func (ar *automator) stateChangeHandler(event evt.Message) {
	if "event" != event.Type || "state_changed" != event.Event.EventType || event.Event.Data.NewState == nil {
		return
	}

	for _, au := range ar.automations {
		//  debounce events and triggers
		if time.Now().Sub(ar.exectimes[au.GetID()]) < time.Second {
			continue
		}

		au.Lock()

		fired := false
		for _, tr := range au.GetTriggers() {
			switch tr.(type) {
			case StateTrigger:
				if tr.(StateTrigger).GetEntityID() != event.Event.Data.NewState.EntityID {
					continue
				}
				fired = fired || tr.(StateTrigger).FireState(*event.Event)
			}
			if fired {
				break
			}
		}

		if fired {
			ar.executeIfConditions(au)
		}

		au.Unlock()
	}
}

func (ar *automator) tickerHandler(now time.Time) {
	for _, au := range ar.automations {
		au.Lock()

		fired := false
		for _, tr := range au.GetTriggers() {
			switch tr.(type) {
			case TimeTrigger:
				fired = fired || tr.(TimeTrigger).FireTime(now)
			}
		}

		if fired {
			ar.executeIfConditions(au)
		}

		au.Unlock()
	}
}

var automatorInstanve Automator
var automatorOnce sync.Once

// NewAutomator ...
func NewAutomator(eventBus evt.Bus) Automator {
	automatorOnce.Do(func() {
		a := &automator{
			eventBus:    eventBus,
			automations: make(map[string]Automation, 0),
			components:  make(map[string]cmp.Component, 0),
			exectimes:   make(map[string]time.Time, 0),
		}

		a.eventBus.Subscribe(RegisterTopic, "Automator.registerAutomation", a.registerAutomation)
		a.eventBus.Subscribe(cmp.RegisterTopic, "Automator.registerComponent", a.registerComponent)
		a.eventBus.SubscribeAsync(evt.StateChangedTopic, "Automator.stateChangeHandler", a.stateChangeHandler, true)
		a.eventBus.SubscribeAsync(cmp.TickerTopic, "Automator.tickerHandler", a.tickerHandler, true)

		automatorInstanve = a
	})

	return automatorInstanve
}
