package user_queue

import (
	"app/internal/marshal"
	"app/internal/mill"
	"app/user"
	"context"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/hibiken/asynq"
)

type PermissionsCleanupQueue struct {
	client *asynq.Client
	codec  marshal.Codec
}

func NewPermissionsCleanupQueue(client *asynq.Client, codec marshal.Codec) *PermissionsCleanupQueue {
	return &PermissionsCleanupQueue{client: client, codec: codec}
}

func (queue *PermissionsCleanupQueue) Enqueue(user user.User) error {
	b, marshalErr := queue.codec.Marshal(user)
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
