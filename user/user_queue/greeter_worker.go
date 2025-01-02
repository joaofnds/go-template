package user_queue

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
)

type GreeterWorker struct{}

func NewGreeterWorker() *GreeterWorker { return &GreeterWorker{} }

func (greeter *GreeterWorker) RegisterQueueHandler(mux *asynq.ServeMux) {
	mux.Handle(TaskGreet, greeter)
}

func (greeter *GreeterWorker) ProcessTask(_ context.Context, task *asynq.Task) error {
	fmt.Printf("[%s] A new user just signed up, welcome %s!\n", task.ResultWriter().TaskID(), task.Payload())
	return nil
}
