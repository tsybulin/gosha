package ptf

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/cmp/intr"
	"github.com/tsybulin/gosha/evt"
	"github.com/tsybulin/gosha/logger"
)

type weather struct {
	Platform
	url        string
	period     time.Duration
	eventBus   evt.Bus
	components map[string]cmp.WeatherSensor
}

func (wp *weather) registerWeatherComponent(c cmp.Component) {
	if c.GetPlatform() != "weather" {
		return
	}

	// WeatherSensor
	wc, ok := c.(cmp.WeatherSensor)
	if !ok {
		return
	}

	wp.components[c.GetID()] = wc

}

func (wp *weather) doget() {
	resp, err := http.Get(wp.url)
	if err != nil {
		if wp.eventBus != nil {
			wp.eventBus.Publish(logger.Topic, logger.LevelError, "WeatherPlatform.doget error", err.Error())
		}
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if wp.eventBus != nil {
			wp.eventBus.Publish(logger.Topic, logger.LevelError, "WeatherPlatform.doget error", err.Error())
		}
		return
	}

	var wi cmp.WeatherInfo
	err = json.Unmarshal(body, &wi)
	if err != nil {
		if wp.eventBus != nil {
			wp.eventBus.Publish(logger.Topic, logger.LevelError, "WeatherPlatform.doget error", err.Error())
		}
		return
	}

	for _, wsold := range wp.components {
		event := eventFor(wsold, func() {})

		wsnew := intr.NewWeather(map[string]string{"weather": wsold.GetID()})

		event.Event.Data.NewState = eventFor(wsnew, func() {
			wsnew.SetAttributes(wi)
		}).Event.Data.NewState

		wp.eventBus.Publish(evt.StateChangedTopic, event)
	}

}

func (wp *weather) Start(eventBus evt.Bus) {
	wp.eventBus = eventBus
	wp.eventBus.Subscribe(cmp.RegisterTopic, "WeatherPlatform.registerWeatherComponent", wp.registerWeatherComponent)
	wp.eventBus.Publish(ReadyTopic, cmp.DomainWeather)

	go func(wp *weather) {
		wp.doget()
	}(wp)

	go func(wp *weather) {
		ticker := time.NewTicker(wp.period)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				wp.doget()
			}
		}

	}(wp)
}

// NewWeatherPlatform ...
func NewWeatherPlatform(cfg map[string]string) Platform {
	p, err := time.ParseDuration(cfg["period"])
	if err != nil {
		p = 30 * time.Minute
	}

	return &weather{
		Platform:   newPlatform("weather"),
		url:        cfg["weather"],
		period:     p,
		components: make(map[string]cmp.WeatherSensor, 0),
	}
}
