package event

type Event[T any] struct {
	cap       int
	listeners []chan T
}

func NewEvent[T any](cap int) *Event[T] {
	return &Event[T]{cap: cap}
}

func (e *Event[T]) Close() {
	for _, c := range e.listeners {
		close(c)
	}
}

func (e *Event[T]) Listen() <-chan T {
	c := make(chan T, e.cap)
	e.listeners = append(e.listeners, c)
	return c
}

func (e *Event[T]) Send(t T) {
	go func() {
		defer func() { _ = recover() }()

		for _, c := range e.listeners {
			c <- t
		}
	}()
}
