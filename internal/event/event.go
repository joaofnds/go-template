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

func On[T any](f Handler[T]) {
	name := fmt.Sprintf("%T", *new(T))

	lock.Lock()
	defer lock.Unlock()

	listeners[name] = append(listeners[name], wrap(f))
}

func Send[T any](t T) {
	name := fmt.Sprintf("%T", t)

	lock.RLock()
	defer lock.RUnlock()

	for _, l := range listeners[name] {
		go l(t)
	}
}

func wrap[T any](f Handler[T]) Handler[any] {
	return func(t any) { f(t.(T)) }
}
