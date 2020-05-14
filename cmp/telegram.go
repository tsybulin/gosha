package cmp

// Telegram ...
type Telegram interface {
	Component
	GetChat() string
	GetBot() string
	GetSubject() string
	GetMessage() string
}
