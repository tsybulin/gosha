package ptf

import (
	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/evt"
	"github.com/tsybulin/gosha/logger"
)

// Manager ...
type Manager interface{}

type manager struct {
	platforms map[string]Platform
	eventBus  evt.Bus
}

func (pm *manager) registerPlatform(p Platform) {
	pm.platforms[p.GetPlatform()] = p
	go func() {
		p.Start(pm.eventBus)
		pm.eventBus.Publish(logger.Topic, logger.LevelInfo, "PlatformManager.registerPlatform", p.GetPlatform())
	}()
}

func (pm *manager) platformPush(c cmp.Component, what string) {
	pid := c.GetPlatform()
	if len(pid) == 0 {
		return
	}

	platform := pm.platforms[pid]

	if platform == nil {
		return
	}

	platform.Push(c, what)

	go func() {
		pm.eventBus.Publish(logger.Topic, logger.LevelDebug, "PlatformManager.platformPush", c.GetID())
	}()
}

// NewManager ...
func NewManager(eventBus evt.Bus) Manager {
	pm := &manager{
		platforms: make(map[string]Platform, 0),
		eventBus:  eventBus,
	}

	eventBus.Subscribe(RegisterTopic, "PlatformManager.registerPlatform", pm.registerPlatform)
	eventBus.SubscribeAsync(evt.PtfPushTopic, "PlatformManager.platformPush", pm.platformPush, false)

	eventBus.Publish(logger.Topic, logger.LevelInfo, "PlatformManager created")

	return pm
}
