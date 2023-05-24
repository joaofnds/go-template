package user

import "context"

type Probe interface {
	FailedToCreateUser(error)
	FailedToDeleteAll(error)
	FailedToFindByName(error)
	FailedToRemoveUser(error, User)
	FailedToEnqueue(error)
	UserCreated(context.Context, User)
}
