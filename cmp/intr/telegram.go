package intr

import (
	"github.com/tsybulin/gosha/cmp"
)

type telegram struct {
	cmp.Component
	bot     string
	chat    string
	subject string
	message string
}

func (t *telegram) GetChat() string {
	return t.chat
}

func (t *telegram) GetBot() string {
	return t.bot
}

func (t *telegram) GetSubject() string {
	return t.subject
}

func (t *telegram) GetMessage() string {
	return t.message
}

// NewTelegram ...
func NewTelegram(cfg map[string]string) cmp.Telegram {
	return &telegram{
		Component: NewComponent(cmp.DomainTelegram, cfg["telegram"], "internal"),
		bot:       cfg["bot"],
		chat:      cfg["chat"],
		subject:   cfg["subject"],
		message:   cfg["message"],
	}
}
