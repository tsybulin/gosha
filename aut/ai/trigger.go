package ai

import "github.com/tsybulin/gosha/aut"

type trigger struct {
	platform string
}

func (t *trigger) GetPlatform() string {
	return t.platform
}

func (*trigger) Fire() bool {
	return false
}

func newTrigger(platform string) aut.Trigger {
	return &trigger{
		platform: platform,
	}
}
