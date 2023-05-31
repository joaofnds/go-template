package event

import (
	"fmt"
	"sync"
)

var (
	lock      = sync.RWMutex{}
	listeners = make(map[string][]Handler[any])
)

type Handler[T any] func(t T)

func On[T any](handler Handler[T]) {
	name := fmt.Sprintf("%T", *new(T))

	lock.Lock()
	defer lock.Unlock()

	listeners[name] = append(listeners[name], wrap(handler))
}

func Emit[T any](value T) {
	name := fmt.Sprintf("%T", value)

	lock.RLock()
	defer lock.RUnlock()

	for _, l := range listeners[name] {
		go l(value)
	}
}

func wrap[T any](handler Handler[T]) Handler[any] {
	return func(value any) { handler(value.(T)) }
}
