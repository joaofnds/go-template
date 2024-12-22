package user_queue

import (
	"app/internal/mill"
	"app/user"
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/hibiken/asynq"
)

type Greeter struct {
	client *asynq.Client
}

func NewGreeter(
	client *asynq.Client,
) *Greeter {
	return &Greeter{
		client: client,
	}
}

func (greeter *Greeter) RegisterEventHandlers(processor *cqrs.EventProcessor) error {
	return processor.AddHandlers(
		mill.NewEventHandler(func(_ context.Context, event *user.UserCreated) error {
			return greeter.Enqueue(event.User.Email)
		}),
	)
}

func (greeter *Greeter) Type() string {
	return "greet"
}

func (greeter *Greeter) RegisterQueueHandler(mux *asynq.ServeMux) {
	mux.Handle(greeter.Type(), greeter)
}

func (greeter *Greeter) Enqueue(email string) error {
	task := asynq.NewTask(greeter.Type(), []byte(email))
	_, err := greeter.client.Enqueue(task)
	return err
}

func (greeter *Greeter) ProcessTask(_ context.Context, task *asynq.Task) error {
	fmt.Printf("[%s] A new user just signed up, welcome %s!\n", task.ResultWriter().TaskID(), task.Payload())
	return nil
}
