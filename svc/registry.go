package svc

import (
	"strconv"
	"sync"

	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/evt"
	"github.com/tsybulin/gosha/logger"
)

// Registry ...
type Registry interface {
	GetDomainService(cmp.Domain) Service
	GetService(string) Service
	Execute(string, string, string)
	Services() []Description
	States() StateResult
}

type registry struct {
	eventBus evt.Bus
	services map[cmp.Domain]Service
}

func (r *registry) GetDomainService(domain cmp.Domain) Service {
	return r.services[domain]
}

func (r *registry) GetService(id string) Service {
	for _, v := range r.services {
		if v.GetID() == id {
			return v
		}
	}

	return nil
}

func (r *registry) States() StateResult {
	stateResults := make([]evt.State, 0)

	for _, s := range r.services {
		ss := s.States()
		if ss.Result != nil {
			for _, sr := range ss.Result {
				stateResults = append(stateResults, sr)
			}
		}
	}

	return StateResult{
		ID:      1,
		Type:    "result",
		Success: true,
		Result:  stateResults,
	}
}

func (r *registry) Execute(service, action, component string) {
	s := r.GetService(service)
	if s == nil {
		r.eventBus.Publish(logger.Topic, logger.LevelWarn, "ServiceRegistry.Execute unknown service:", service)
		return
	}

	switch s.(type) {
	case Light:
		switch action {
		case "turn_on":
			s.(Light).TurnOn(component)
		case "turn_off":
			s.(Light).TurnOff(component)
		case "toggle":
			s.(Light).Toggle(component)
		case "10", "20", "30", "40", "50", "60", "70", "80", "90", "100":
			if b, err := strconv.Atoi(action); err == nil {
				s.(Light).SetBrightness(component, int16(b))
			}
		default:
			if b, err := strconv.Atoi(action); err == nil {
				s.(Light).SetBrightness(component, int16(b))
			} else {
				r.eventBus.Publish(logger.Topic, logger.LevelWarn, "ServiceRegistry.Execute service", service, " unsupported action", action)
			}
		}
	case Switch:
		switch action {
		case "turn_on":
			s.(Switch).TurnOn(component)
		case "turn_off":
			s.(Switch).TurnOff(component)
		case "toggle":
			s.(Switch).Toggle(component)
		default:
			r.eventBus.Publish(logger.Topic, logger.LevelWarn, "ServiceRegistry.Execute service", service, " unsupported action", action)
		}
	case Group:
		switch action {
		case "turn_on":
			s.(Group).TurnOn(component)
		case "turn_off":
			s.(Group).TurnOff(component)
		case "toggle":
			s.(Group).Toggle(component)
		default:
			r.eventBus.Publish(logger.Topic, logger.LevelWarn, "ServiceRegistry.Execute service", service, " unsupported action", action)
		}
	case Timer:
		switch action {
		case "start":
			s.(Timer).Start(component)
		case "stop":
			s.(Timer).Stop(component)
		case "cancel":
			s.(Timer).Cancel(component)
		default:
			r.eventBus.Publish(logger.Topic, logger.LevelWarn, "ServiceRegistry.Execute service", service, " unsupported action", action)
		}
	case Script:
		switch action {
		case "execute":
			s.(Script).Execute(component)
		default:
			r.eventBus.Publish(logger.Topic, logger.LevelWarn, "ServiceRegistry.Execute service", service, " unsupported action", action)
		}
	case Telegram:
		switch action {
		case "notify":
			s.(Telegram).Notify(component)
		default:
			r.eventBus.Publish(logger.Topic, logger.LevelWarn, "ServiceRegistry.Execute service", service, " unsupported action", action)
		}
	case Alarm:
		switch action {
		case "disarm":
			s.(Alarm).Disarm(component)
		case "arm_home":
			s.(Alarm).ArmHome(component)
		case "arm_away":
			s.(Alarm).ArmAway(component)
		case "arm_night":
			s.(Alarm).ArmNight(component)
		case "trigger":
			s.(Alarm).Trigger(component)
		default:
			r.eventBus.Publish(logger.Topic, logger.LevelWarn, "ServiceRegistry.Execute service", service, " unsupported action", action)
		}
	case Countdown:
		switch action {
		case "start":
			s.(Countdown).Start(component)
		case "stop":
			s.(Countdown).Stop(component)
		default:
			r.eventBus.Publish(logger.Topic, logger.LevelWarn, "ServiceRegistry.Execute service", service, " unsupported action", action)
		}
	default:
		r.eventBus.Publish(logger.Topic, logger.LevelWarn, "ServiceRegistry.Execute unsupported service type for:", service)
	}
}

func (r *registry) Services() []Description {
	svcs := make([]Description, 0)
	for _, s := range r.services {
		d := s.Description()
		if d.Methods != nil {
			svcs = append(svcs, s.Description())
		}
	}
	return svcs
}

var registryInstance Registry
var registryOnce sync.Once

// NewRegistry ...
func NewRegistry(eventBus evt.Bus) Registry {

	registryOnce.Do(func() {
		r := &registry{
			eventBus: eventBus,
			services: make(map[cmp.Domain]Service, 0),
		}

		r.services[cmp.DomainGroup] = newGroupService(eventBus)
		r.services[cmp.DomainLight] = newLightService(eventBus)
		r.services[cmp.DomainMqtt] = newMqttService(eventBus)
		r.services[cmp.DomainBinarySensor] = newBinarySensorService(eventBus)
		r.services[cmp.DomainSensor] = newSensorService(eventBus)
		r.services[cmp.DomainSwitch] = newSwitchService(eventBus)
		r.services[cmp.DomainTimer] = newTimerService(eventBus)
		r.services[cmp.DomainScript] = newScriptService(eventBus)
		r.services[cmp.DomainWeather] = newSensorWeather(eventBus)
		r.services[cmp.DomainTelegram] = newTelegram(eventBus)
		r.services[cmp.DomainAlarm] = newAlarmService(eventBus)
		r.services[cmp.DomainCountdown] = newCountdownService(eventBus)

		registryInstance = r
	})

	return registryInstance
}
