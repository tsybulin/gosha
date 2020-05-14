package cmp

import (
	"time"

	"github.com/tsybulin/gosha/evt"
)

const TickerTopic = "ticker:tick"

type Ticker struct {
	eventBus evt.Bus
	ticker   *time.Ticker
}

func (t *Ticker) Start() {
	t.ticker = time.NewTicker(time.Second)
	go func(t *Ticker) {
		for now := range t.ticker.C {
			t.eventBus.Publish(TickerTopic, now)
		}
	}(t)
}

func (t *Ticker) Stop() {
	if t.ticker != nil {
		t.ticker.Stop()
		t.ticker = nil
	}
}

func NewTicker(eventBus evt.Bus) *Ticker {
	t := &Ticker{
		eventBus: eventBus,
	}

	return t
}
