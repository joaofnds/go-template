package user

type UserCreated struct {
	User User
}

type FailedToCreateUser struct {
	User User
	Err  error
}

type FailedToDeleteAll struct {
	Err error
}

type FailedToFindByName struct {
	Err error
}

type FailedToRemoveUser struct {
	User User
	Err  error
}
