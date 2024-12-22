package user

type UserCreated struct {
	User User
}

type FailedToCreateUser struct {
	Err error
}

type FailedToDeleteAll struct {
	Err error
}

type FailedToFindByID struct {
	Err error
}

type FailedToFindByName struct {
	Err error
}

type FailedToRemoveUser struct {
	User User
	Err  error
}
