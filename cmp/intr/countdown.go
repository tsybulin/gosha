package intr

import (
	"fmt"
	"time"

	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/evt"
)

type countdown struct {
	cmp.Component
	duration time.Duration
	left     time.Duration
}

func (c *countdown) GetDuration() time.Duration {
	return c.duration
}

func (c *countdown) SetDuration(duration time.Duration) {
	c.duration = duration
}

func (c *countdown) GetLeft() time.Duration {
	return c.left
}

func (c *countdown) SetLeft(left time.Duration) {
	c.left = left
}

func (c *countdown) GetState() evt.State {
	state := c.Component.GetState()
	state.Attributes["duration"] = c.duration
	state.Attributes["left"] = c.left
	state.State = fmt.Sprint(c.left.Minutes())
	return state
}

// NewCountdown ...
func NewCountdown(cfg map[string]string) cmp.Countdown {
	d, err := time.ParseDuration(cfg["duration"])
	if err != nil {
		return nil
	}

	return &countdown{
		Component: NewComponent(cmp.DomainCountdown, cfg["countdown"], "internal"),
		duration:  d,
		left:      0,
	}
}
