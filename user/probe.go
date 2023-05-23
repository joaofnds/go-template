package user

type Probe interface {
	FailedToCreateUser(error)
	FailedToDeleteAll(error)
	FailedToFindByName(error)
	FailedToRemoveUser(error, User)
	FailedToEnqueue(error)
	UserCreated()
}
