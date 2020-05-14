package svc

import (
	"time"

	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/evt"
	"github.com/tsybulin/gosha/logger"
)

// Timer ...
type Timer interface {
	Service
	Start(string)
	Stop(string)
	Cancel(string)
}

type timerService struct {
	Service
	eventBus evt.Bus
	timers   map[string]cmp.Timer
}

func (ts *timerService) Components() []string {
	components := make([]string, 0)
	for id := range ts.timers {
		components = append(components, id)
	}
	return components
}

func (ts *timerService) States() StateResult {
	stateResults := make([]evt.State, 0)

	for _, c := range ts.timers {
		stateResults = append(stateResults, c.GetState())
	}

	return StateResult{
		ID:      1,
		Type:    "result",
		Success: true,
		Result:  stateResults,
	}
}

func (ts *timerService) Description() Description {
	return Description{
		ID: ts.GetID(),
		Methods: []Method{
			{Method: "start", Parameters: []string{"timer_id"}},
			{Method: "stop", Parameters: []string{"timer_id"}},
			{Method: "cancel", Parameters: []string{"timer_id"}},
		},
	}
}

func (ts *timerService) tick(time.Time) {
	for _, timer := range ts.timers {
		if timer.GetTimerState() == cmp.TimerStateInactive {
			continue
		}

		timer.SetLeft(timer.GetLeft() - time.Second)
		if timer.GetLeft() <= 0 {
			ts.setState("tick", timer.GetID(), cmp.TimerStateInactive, true)
		} else {
			go func() {
				ts.eventBus.Publish(logger.Topic, logger.LevelDebug, "TimerService.tick:", timer.GetID(), timer.GetLeft())
			}()
		}

	}
}

func (ts *timerService) registerTimer(c cmp.Component) {
	if c.GetDomain() != cmp.DomainTimer {
		return
	}

	t, ok := c.(cmp.Timer)

	if !ok {
		return
	}

	ts.timers[c.GetID()] = t

	go func() {
		ts.eventBus.Publish(logger.Topic, logger.LevelDebug, "TimerService.registerTimer", c.GetID())
	}()
}

func (ts *timerService) setState(op, id string, state cmp.TimerState, notify bool) {
	t := ts.timers[id]

	if t == nil {
		go func() {
			ts.eventBus.Publish(logger.Topic, logger.LevelWarn, "TimerService.setState unknown id:", id)
		}()
		return
	}

	if t.GetTimerState() == state {
		return
	}

	event := eventFor(t, func() {
		if state == cmp.TimerStateActive {
			t.SetLeft(t.GetDuration())
		} else {
			t.SetLeft(0)
		}
		t.SetTimerState(state)
	})

	go func() {
		if notify {
			ts.eventBus.Publish(evt.StateChangedTopic, event)
		}

		ts.eventBus.Publish(logger.Topic, logger.LevelSystem, "TimerService.", op, id, state)
	}()
}

// Start ...
func (ts *timerService) Start(id string) {
	ts.setState("start", id, cmp.TimerStateActive, true)
}

// Stop ...
func (ts *timerService) Stop(id string) {
	ts.setState("stop", id, cmp.TimerStateInactive, true)
}

// Cancel ...
func (ts *timerService) Cancel(id string) {
	ts.setState("cancel", id, cmp.TimerStateInactive, false)
}

func newTimerService(eventBus evt.Bus) Timer {
	ts := &timerService{
		Service:  newService(cmp.DomainTimer),
		eventBus: eventBus,
		timers:   make(map[string]cmp.Timer, 0),
	}

	ts.eventBus.SubscribeAsync(cmp.TickerTopic, "TimerService.tick", ts.tick, true)
	ts.eventBus.Subscribe(cmp.RegisterTopic, "TimerService.registerTimer", ts.registerTimer)

	return ts
}
