package evt

import (
	"reflect"
	"sync"
)

// Bus ...
type Bus interface {
	Subscribe(topic, id string, fn interface{})
	SubscribeAsync(topic, id string, fn interface{}, transactional bool)
	SubscribeOnce(topic, id string, fn interface{})
	SubscribeOnceAsync(topic, id string, fn interface{})
	HasCallback(topic string) bool
	Unsubscribe(topic, id string)
	Publish(topic string, args ...interface{})
	WaitAsync()
}

type eventHandler struct {
	callBack      reflect.Value
	flagOnce      bool
	async         bool
	transactional bool
	sync.Mutex
}

type eventBus struct {
	handlers map[string]map[string]*eventHandler
	lock     sync.Mutex
	wg       sync.WaitGroup
}

func (bus *eventBus) subscribe(topic, id string, fn interface{}, handler *eventHandler) {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	if _, ok := bus.handlers[topic]; !ok {
		bus.handlers[topic] = make(map[string]*eventHandler)
	}
	bus.handlers[topic][id] = handler
}

func (bus *eventBus) Subscribe(topic, id string, fn interface{}) {
	bus.subscribe(topic, id, fn, &eventHandler{
		reflect.ValueOf(fn),
		false,
		false,
		false,
		sync.Mutex{},
	})
}

func (bus *eventBus) SubscribeAsync(topic, id string, fn interface{}, transactional bool) {
	bus.subscribe(topic, id, fn, &eventHandler{
		reflect.ValueOf(fn),
		false,
		true,
		transactional,
		sync.Mutex{},
	})
}

func (bus *eventBus) SubscribeOnce(topic, id string, fn interface{}) {
	bus.subscribe(topic, id, fn, &eventHandler{
		reflect.ValueOf(fn),
		true,
		false,
		false,
		sync.Mutex{},
	})
}

func (bus *eventBus) SubscribeOnceAsync(topic, id string, fn interface{}) {
	bus.subscribe(topic, id, fn, &eventHandler{
		reflect.ValueOf(fn),
		true,
		true,
		false,
		sync.Mutex{},
	})
}

func (bus *eventBus) HasCallback(topic string) bool {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	_, ok := bus.handlers[topic]
	if ok {
		return len(bus.handlers[topic]) > 0
	}
	return false
}

func (bus *eventBus) unsubscribe(topic, id string) {
	if _, ok := bus.handlers[topic][id]; ok {
		delete(bus.handlers[topic], id)
	}
}

func (bus *eventBus) Unsubscribe(topic, id string) {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	bus.unsubscribe(topic, id)
}

func (bus *eventBus) setUpPublish(callback *eventHandler, args ...interface{}) []reflect.Value {
	funcType := callback.callBack.Type()
	passedArguments := make([]reflect.Value, len(args))
	for i, v := range args {
		if v == nil {
			passedArguments[i] = reflect.New(funcType.In(i)).Elem()
		} else {
			passedArguments[i] = reflect.ValueOf(v)
		}
	}

	return passedArguments
}

func (bus *eventBus) publish(handler *eventHandler, args ...interface{}) {
	passedArguments := bus.setUpPublish(handler, args...)
	handler.callBack.Call(passedArguments)
}

func (bus *eventBus) publishAsync(handler *eventHandler, args ...interface{}) {
	defer bus.wg.Done()
	if handler.transactional {
		defer handler.Unlock()
	}
	bus.publish(handler, args...)
}

func (bus *eventBus) Publish(topic string, args ...interface{}) {
	bus.lock.Lock()
	defer bus.lock.Unlock()

	remove := make([]string, 0)

	for id, handler := range bus.handlers[topic] {
		if handler.flagOnce {
			remove = append(remove, id)
		}
		if !handler.async {
			bus.publish(handler, args...)
		} else {
			bus.wg.Add(1)
			if handler.transactional {
				bus.lock.Unlock()
				handler.Lock()
				bus.lock.Lock()
			}
			go bus.publishAsync(handler, args...)
		}
	}

	for _, id := range remove {
		bus.unsubscribe(topic, id)
	}
}

func (bus *eventBus) WaitAsync() {
	bus.wg.Wait()
}

// NewBus ...
func NewBus() Bus {
	b := &eventBus{
		make(map[string]map[string]*eventHandler),
		sync.Mutex{},
		sync.WaitGroup{},
	}
	return b
}
