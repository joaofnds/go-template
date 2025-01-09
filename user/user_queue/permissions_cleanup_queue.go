package user_queue

import (
	"app/internal/mill"
	"app/user"
	"context"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/bytedance/sonic"
	"github.com/hibiken/asynq"
)

type PermissionsCleanupQueue struct {
	client *asynq.Client
}

func NewPermissionsCleanupQueue(client *asynq.Client) *PermissionsCleanupQueue {
	return &PermissionsCleanupQueue{client: client}
}

func (queue *PermissionsCleanupQueue) Enqueue(user user.User) error {
	b, marshalErr := sonic.Marshal(user)
	if marshalErr != nil {
		return marshalErr
	}

	task := asynq.NewTask(TaskPermissionCleanup, b)
	_, err := queue.client.Enqueue(task)
	return err
}

func (queue *PermissionsCleanupQueue) RegisterEventHandlers(processor *cqrs.EventProcessor) error {
	return processor.AddHandlers(
		mill.NewEventHandler(func(_ context.Context, event *user.UserRemoved) error {
			return queue.Enqueue(event.User)
		}),
	)
}
