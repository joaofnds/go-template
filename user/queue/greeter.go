package queue

import (
	"app/internal/event"
	"app/user"
	"context"
	"fmt"

	"github.com/hibiken/asynq"
)

type Greeter struct {
	client *asynq.Client
}

func NewGreeter(client *asynq.Client) *Greeter {
	return &Greeter{client: client}
}

func (g *Greeter) Listen() {
	event.On(func(e user.UserCreated) { _ = g.Enqueue(e.User.Name) })
}

func (g *Greeter) Type() string {
	return "greet"
}

func (g *Greeter) Register(mux *asynq.ServeMux) {
	mux.Handle(g.Type(), g)
}

func (g *Greeter) Enqueue(userName string) error {
	task := asynq.NewTask(g.Type(), []byte(userName))
	_, err := g.client.Enqueue(task)
	return err
}

func (g *Greeter) ProcessTask(_ context.Context, task *asynq.Task) error {
	fmt.Printf("[%s] A new user just signed up, welcome %s!\n", task.ResultWriter().TaskID(), task.Payload())
	return nil
}
