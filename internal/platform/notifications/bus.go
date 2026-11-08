package notifications

import "sync"

type EventBus struct {
	listeners map[string][]any
	mu        sync.RWMutex
}

var defaultBus = &EventBus{listeners: make(map[string][]any)}

func On(eventName string, fn any) {
	defaultBus.mu.Lock()
	defer defaultBus.mu.Unlock()
	defaultBus.listeners[eventName] = append(defaultBus.listeners[eventName], fn)
}

func Dispatch(eventName string, data any) {
	defaultBus.mu.RLock()
	defer defaultBus.mu.RUnlock()
	for _, fn := range defaultBus.listeners[eventName] {
		if handler, ok := fn.(func(any)); ok {
			go handler(data)
		}
	}
}
