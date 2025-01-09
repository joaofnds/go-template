package user_queue

import (
	"app/authz"
	"app/user"
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/hibiken/asynq"
)

type PermissionsCleanupWorker struct {
	roles       authz.RoleManager
	permissions authz.PermissionManager
}

func NewPermissionsCleanupWorker(
	roles authz.RoleManager,
	permissions authz.PermissionManager,
) *PermissionsCleanupWorker {
	return &PermissionsCleanupWorker{
		roles:       roles,
		permissions: permissions,
	}
}

func (worker *PermissionsCleanupWorker) RegisterQueueHandler(mux *asynq.ServeMux) {
	mux.Handle(TaskPermissionCleanup, worker)
}

func (worker *PermissionsCleanupWorker) ProcessTask(_ context.Context, task *asynq.Task) error {
	var u user.User
	if err := sonic.Unmarshal(task.Payload(), &u); err != nil {
		return err
	}

	if err := worker.roles.RevokeAll(user.NewRef(u.ID)); err != nil {
		return fmt.Errorf("failed to revoke roles: %w", err)
	}

	if err := worker.permissions.RemoveBySubject(user.NewRef(u.ID)); err != nil {
		return fmt.Errorf("failed to remove permissions: %w", err)
	}

	return nil
}
