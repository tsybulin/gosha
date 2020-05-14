package logger

import (
	"github.com/tsybulin/gosha/evt"
)

type EventLogger struct {
	eventBus evt.Bus
}

func (el *EventLogger) log(event evt.Message) {
	go func() {
		el.eventBus.Publish(Topic, LevelInfo, "EventLogger", event.String())
	}()
}

func (el *EventLogger) Start() {
	el.eventBus.Subscribe(evt.StateChangedTopic, "EventLogger.log", el.log)

}

func NewEventLogger(eventBus evt.Bus) *EventLogger {
	el := &EventLogger{
		eventBus: eventBus,
	}
	return el
}
