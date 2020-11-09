package svc

import (
	"time"

	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/evt"
	"github.com/tsybulin/gosha/logger"
)

// Countdown service
type Countdown interface {
	Service
	Start(string)
	Stop(string)
}

type countdownService struct {
	Service
	eventBus evt.Bus
	cntdowns map[string]cmp.Countdown
}

func (cs *countdownService) Components() []string {
	components := make([]string, 0)
	for id := range cs.cntdowns {
		components = append(components, id)
	}
	return components
}

func (cs *countdownService) States() StateResult {
	stateResults := make([]evt.State, 0)

	for _, c := range cs.cntdowns {
		stateResults = append(stateResults, c.GetState())
	}

	return StateResult{
		ID:      1,
		Type:    "result",
		Success: true,
		Result:  stateResults,
	}
}

func (cs *countdownService) Description() Description {
	return Description{
		ID: cs.GetID(),
		Methods: []Method{
			{Method: "start", Parameters: []string{"cntdwn_id"}},
			{Method: "stop", Parameters: []string{"cntdwn_id"}},
		},
	}
}

func (cs *countdownService) Start(id string) {
	c := cs.cntdowns[id]

	if c == nil {
		go func() {
			cs.eventBus.Publish(logger.Topic, logger.LevelWarn, "CountdownService.Start unknown id:", id)
		}()
		return
	}

	cs.setLeft(c, c.GetDuration())
}

func (cs *countdownService) Stop(id string) {
	c := cs.cntdowns[id]

	if c == nil {
		go func() {
			cs.eventBus.Publish(logger.Topic, logger.LevelWarn, "CountdownService.Stop unknown id:", id)
		}()
		return
	}

	cs.setLeft(c, 0)
}

func (cs *countdownService) setLeft(c cmp.Countdown, left time.Duration) {
	event := eventFor(c, func() {
		c.SetLeft(left)
	})
	go func() {
		cs.eventBus.Publish(evt.StateChangedTopic, event)
		cs.eventBus.Publish(logger.Topic, logger.LevelDebug, "CountdownService.setSLeft:", c.GetID(), c.GetLeft())
	}()
}

func (cs *countdownService) tick(time.Time) {
	now := time.Now()

	if now.Second() != 0 {
		return
	}

	for _, c := range cs.cntdowns {
		if c.GetLeft() <= 0 {
			continue
		}

		cs.setLeft(c, c.GetLeft()-time.Minute)
	}
}

func (cs *countdownService) registerCountdown(c cmp.Component) {
	if c.GetDomain() != cmp.DomainCountdown {
		return
	}

	t, ok := c.(cmp.Countdown)

	if !ok {
		return
	}

	cs.cntdowns[c.GetID()] = t

	go func() {
		cs.eventBus.Publish(logger.Topic, logger.LevelDebug, "CountdownService.registerCountdown", c.GetID())
	}()
}

func newCountdownService(eventBus evt.Bus) Countdown {
	cs := &countdownService{
		Service:  newService(cmp.DomainCountdown),
		eventBus: eventBus,
		cntdowns: make(map[string]cmp.Countdown, 0),
	}

	cs.eventBus.SubscribeAsync(cmp.TickerTopic, "CountdownService.tick", cs.tick, true)
	cs.eventBus.Subscribe(cmp.RegisterTopic, "CountdownService.registerCountdown", cs.registerCountdown)

	return cs
}
