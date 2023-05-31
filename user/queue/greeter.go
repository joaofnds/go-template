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

func (greeter *Greeter) Listen() {
	event.On(func(e user.UserCreated) { _ = greeter.Enqueue(e.User.Name) })
}

func (greeter *Greeter) Type() string {
	return "greet"
}

func (greeter *Greeter) Register(mux *asynq.ServeMux) {
	mux.Handle(greeter.Type(), greeter)
}

func (greeter *Greeter) Enqueue(userName string) error {
	task := asynq.NewTask(greeter.Type(), []byte(userName))
	_, err := greeter.client.Enqueue(task)
	return err
}

func (greeter *Greeter) ProcessTask(_ context.Context, task *asynq.Task) error {
	fmt.Printf("[%s] A new user just signed up, welcome %s!\n", task.ResultWriter().TaskID(), task.Payload())
	return nil
}
