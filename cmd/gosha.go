package main

import (
	"context"

	"github.com/tsybulin/gosha"
	"github.com/tsybulin/gosha/aut"
	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/evt"
	"github.com/tsybulin/gosha/logger"
	"github.com/tsybulin/gosha/ptf"
	"github.com/tsybulin/gosha/svc"
)

var (
	thelogger       *logger.Logger
	ticker          *cmp.Ticker
	platformManager ptf.Manager
	serviceregistry svc.Registry
	automator       aut.Automator
	configurator    *gosha.Configurator
)

func main() {

	eventBus := evt.NewBus()

	thelogger = logger.NewLogger(eventBus)
	thelogger.SetSeverity(logger.LevelWarn)
	// thelogger.SetSeverity(logger.LevelDebug)

	eventBus.Publish(logger.Topic, logger.LevelSystem, `
***********************************************
*                                             *
*             G    O    S    H    A           *
*                                             *
***********************************************`)

	el := logger.NewEventLogger(eventBus)
	el.Start()

	ticker = cmp.NewTicker(eventBus)
	ticker.Start()

	platformManager = ptf.NewManager(eventBus)
	serviceregistry = svc.NewRegistry(eventBus)
	automator = aut.NewAutomator(eventBus)

	configurator = gosha.NewConfigurator(eventBus, "./config", "config.yaml")
	configurator.LoadConfig()

	ctx := context.Background()
	<-ctx.Done()
}
