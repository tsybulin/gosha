package gosha

import (
	"io/ioutil"
	"os"
	"sync"

	"github.com/tsybulin/gosha/aut"
	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/evt"
	"github.com/tsybulin/gosha/logger"
	"github.com/tsybulin/gosha/ptf"
	"github.com/tsybulin/gosha/ws"
	"gopkg.in/yaml.v2"
)

// Config ...
type Config struct {
	Includes   []string
	Platforms  []map[string]string
	Components []map[string]string
	Groups     []struct {
		Group    string
		Entities []string
	}
	Automations []struct {
		Automation string
		Triggers   []map[string]string
		Conditions []map[string]string
		Actions    []map[string]string
	}
	Scripts []struct {
		Script  string
		Actions []map[string]string
	}
}

// Configurator ...
type Configurator struct {
	eventBus evt.Bus
	basedir  string
	file     string
}

// NewConfigurator ...
func NewConfigurator(eventBus evt.Bus, basedir, file string) *Configurator {
	cg := &Configurator{
		eventBus: eventBus,
		basedir:  basedir,
		file:     file,
	}

	eventBus.SubscribeAsync(evt.CfgReloadAutomationsTopic, "Configurator.handleReloadAutomations", cg.handleReloadAutomations, false)
	eventBus.Publish(logger.Topic, logger.LevelInfo, "Configurator started")

	return cg
}

func (cg *Configurator) parseConfig(basedir, file string) Config {
	config := Config{}

	path := basedir + string(os.PathSeparator) + file

	yamlFile, err := os.Open(path)
	if err != nil {
		cg.eventBus.Publish(logger.Topic, logger.LevelError, "Configurator.ParseConfig error", err)
		return config
	}

	defer yamlFile.Close()

	bs, err := ioutil.ReadAll(yamlFile)
	if err != nil {
		cg.eventBus.Publish(logger.Topic, logger.LevelError, "Configurator.ParseConfig error", err)
		return config
	}

	err = yaml.Unmarshal(bs, &config)
	if err != nil {
		cg.eventBus.Publish(logger.Topic, logger.LevelError, "Configurator.ParseConfig error", err)
		return config
	}

	if config.Includes != nil {
		for _, include := range config.Includes {
			cfg := cg.parseConfig(basedir, include)

			if cfg.Platforms != nil {
				if config.Platforms == nil {
					config.Platforms = cfg.Platforms
				} else {
					for _, v := range cfg.Platforms {
						config.Platforms = append(config.Platforms, v)
					}
				}
			}

			if cfg.Components != nil {
				if config.Components == nil {
					config.Components = cfg.Components
				} else {
					for _, v := range cfg.Components {
						config.Components = append(config.Components, v)
					}
				}
			}

			if cfg.Groups != nil {
				if config.Groups == nil {
					config.Groups = cfg.Groups
				} else {
					for _, v := range cfg.Groups {
						config.Groups = append(config.Groups, v)
					}
				}
			}

			if cfg.Automations != nil {
				if config.Automations == nil {
					config.Automations = cfg.Automations
				} else {
					for _, v := range cfg.Automations {
						config.Automations = append(config.Automations, v)
					}
				}
			}

			if cfg.Scripts != nil {
				if config.Scripts == nil {
					config.Scripts = cfg.Scripts
				} else {
					for _, v := range cfg.Scripts {
						config.Scripts = append(config.Scripts, v)
					}
				}
			}
		}
	}

	return config
}

func (cg *Configurator) loadAatomations(config Config) {
	if config.Automations != nil {
		for _, a := range NewAutomations(config.Automations) {
			cg.eventBus.Publish(aut.RegisterTopic, a)
		}
	}
}

func (cg *Configurator) handleReloadAutomations(bool) {
	config := cg.parseConfig(cg.basedir, cg.file)
	cg.loadAatomations(config)
}

// LoadConfig ...
func (cg *Configurator) LoadConfig() {
	config := cg.parseConfig(cg.basedir, cg.file)

	cg.loadAatomations(config)

	if config.Scripts != nil {
		for _, s := range NewScripts(config.Scripts) {
			cg.eventBus.Publish(aut.ScriptRegisterTopic, s)
		}
	}

	var pwg sync.WaitGroup

	prh := func(domain cmp.Domain) {
		pwg.Done()
	}

	cg.eventBus.SubscribeAsync(ptf.ReadyTopic, "Configurator.prh", prh, false)

	for _, cfg := range config.Platforms {
		pwg.Add(1)

		if len(cfg["mqtt"]) > 0 {
			p := ptf.NewMqttPlatform(cfg)
			cg.eventBus.Publish(ptf.RegisterTopic, p)
		}

		if len(cfg["weather"]) > 0 {
			p := ptf.NewWeatherPlatform(cfg)
			cg.eventBus.Publish(ptf.RegisterTopic, p)
		}

		if len(cfg["homekit"]) > 0 {
			p := ptf.NewHomeKitPlatform(cfg)
			cg.eventBus.Publish(ptf.RegisterTopic, p)
		}

		if len(cfg["web"]) > 0 {
			ws.NewWebService(cg.eventBus, cfg["web"], cfg["auth_key"]).Start()
			pwg.Done()
		}

	}

	pwg.Wait()
	cg.eventBus.Unsubscribe(ptf.ReadyTopic, "Configurator.prh")

	if config.Groups != nil {
		for _, gc := range config.Groups {
			g := CreateGroup(gc.Group, gc.Entities)
			cg.eventBus.Publish(cmp.RegisterGroupTopic, g)
		}
	}

	for _, cfg := range config.Components {
		if c := CreateComponent(cfg); c != nil {
			cg.eventBus.Publish(cmp.RegisterTopic, c)
		} else {
			cg.eventBus.Publish(logger.Topic, logger.LevelError, "Configurator.LoadConfig cant create component", cfg)
		}
	}

	cg.eventBus.Publish(cmp.RegistrationFinishedTopic, true)

}
