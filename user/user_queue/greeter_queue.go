package user_queue

import (
	"app/internal/mill"
	"app/user"
	"context"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/hibiken/asynq"
)

type GreeterQueue struct {
	client *asynq.Client
}

func NewGreeterQueue(client *asynq.Client) *GreeterQueue {
	return &GreeterQueue{client: client}
}

func (queue *GreeterQueue) Enqueue(email string) error {
	task := asynq.NewTask(TaskGreet, []byte(email))
	_, err := queue.client.Enqueue(task)
	return err
}

func (queue *GreeterQueue) RegisterEventHandlers(processor *cqrs.EventProcessor) error {
	return processor.AddHandlers(
		mill.NewEventHandler(func(_ context.Context, event *user.UserCreated) error {
			return queue.Enqueue(event.User.Email)
		}),
	)
}
