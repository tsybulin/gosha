package ai

import "github.com/tsybulin/gosha/aut"

type condition struct {
	platform string
}

func (c *condition) Satisfied() bool {
	return true
}

func (c *condition) GetPlatform() string {
	return c.platform
}

func newCondition(platform string) aut.Condition {
	return &condition{
		platform: platform,
	}
}
