package intr

import "github.com/tsybulin/gosha/cmp"

type group struct {
	cmp.Component
	ids        map[string]string
	components map[string]cmp.Component
}

// AddComponent ...
func (g *group) AddComponent(c cmp.Component) bool {
	if e := g.ids[c.GetID()]; len(e) == 0 {
		return false
	}

	g.components[c.GetID()] = c
	return true
}

// GetComponent ...
func (g *group) GetComponent(id string) cmp.Component {
	return g.components[id]
}

func (g *group) GetComponents() []cmp.Component {
	cs := make([]cmp.Component, 0)
	for _, c := range g.components {
		cs = append(cs, c)
	}
	return cs
}

// IsOn ...
func (g *group) IsOn() bool {
	state := true
	for _, c := range g.components {
		if e, ok := c.(cmp.Switchable); ok && !e.IsOn() {
			state = false
			break
		}

		if e, ok := c.(cmp.BinarySensor); ok && !e.IsOn() {
			state = false
			break
		}
	}

	return state
}

// NewGroup ...
func NewGroup(id string, es []string) cmp.Group {
	g := &group{
		Component:  NewComponent(cmp.DomainGroup, id, "internal"),
		ids:        make(map[string]string, 0),
		components: make(map[string]cmp.Component, 0),
	}

	for _, e := range es {
		g.ids[e] = e
	}

	return g
}
