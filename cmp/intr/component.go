package intr

import (
	"strings"

	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/evt"
)

type component struct {
	id       string
	platform string
	domain   cmp.Domain
}

func (c *component) GetID() string {
	return c.id
}

func (c *component) GetPlatform() string {
	return c.platform
}

func (c *component) GetDomain() cmp.Domain {
	return c.domain
}

func (c *component) GetState() evt.State {
	return evt.State{
		EntityID:   c.id,
		State:      "",
		Attributes: make(map[string]interface{}, 0),
	}
}

// NewComponent ...
func NewComponent(domain cmp.Domain, id, platform string) cmp.Component {
	return &component{
		domain:   domain,
		platform: platform,
		id:       domain.String() + "." + strings.ToLower(id),
	}
}
