package svc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/evt"
	"github.com/tsybulin/gosha/logger"
)

const urlTemplate = "https://api.telegram.org/bot%v/sendMessage"
const messageTemplate = `{"chat_id":%v,"parse_mode": "Markdown","text":"*%v*\n%v"}`

// Telegram ...
type Telegram interface {
	Service
	Notify(string)
}

type telegram struct {
	Service
	eventBus  evt.Bus
	telegrams map[string]cmp.Telegram
}

func (t *telegram) Components() []string {
	return []string{"telegram"}
}

func (t *telegram) Description() Description {
	return Description{
		ID: t.GetID(),
	}
}

func (t *telegram) States() StateResult {
	return StateResult{
		ID:      1,
		Type:    "result",
		Success: true,
	}
}

func (t *telegram) Notify(id string) {
	go func() {
		if tg := t.telegrams[id]; tg != nil {
			t.eventBus.Publish(logger.Topic, logger.LevelInfo, "Telegram.Notify", id)

			msg := fmt.Sprintf(messageTemplate, tg.GetChat(), tg.GetSubject(), tg.GetMessage())
			url := fmt.Sprintf(urlTemplate, tg.GetBot())
			b := bytes.NewBuffer([]byte(msg))

			if resp, err := http.Post(url, "application/json", b); err == nil {
				if body, err := ioutil.ReadAll(resp.Body); err == nil {
					defer resp.Body.Close()
					t.eventBus.Publish(logger.Topic, logger.LevelDebug, "Telegram.Send", string(body))
				} else {
					t.eventBus.Publish(logger.Topic, logger.LevelError, "Telegram.Send", err.Error())
				}
			} else {
				t.eventBus.Publish(logger.Topic, logger.LevelError, "Telegram.Send", err.Error())
			}

		} else {
			t.eventBus.Publish(logger.Topic, logger.LevelWarn, "Telegram.Notify unknown id:", id)
		}

	}()
}

func (t *telegram) registerTelegram(c cmp.Component) {
	if c.GetDomain() != cmp.DomainTelegram {
		return
	}

	tg, ok := c.(cmp.Telegram)

	if !ok {
		return
	}

	t.telegrams[c.GetID()] = tg

	go func() {
		t.eventBus.Publish(logger.Topic, logger.LevelDebug, "Telegram.register", c.GetID())
	}()
}

func newTelegram(eventBus evt.Bus) Service {
	t := &telegram{
		Service:   newService(cmp.DomainTelegram),
		eventBus:  eventBus,
		telegrams: make(map[string]cmp.Telegram, 0),
	}

	t.eventBus.Subscribe(cmp.RegisterTopic, "Telegram.registerTelegram", t.registerTelegram)

	return t
}
