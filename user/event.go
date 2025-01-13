package user

type UserCreated struct {
	User User
}

type UserRemoved struct {
	User User
}

type FailedToCreateUser struct {
	Error error
}

type FailedToDeleteAll struct {
	Error error
}

type FailedToFindByID struct {
	Error error
}

type FailedToFindByName struct {
	Error error
}

type FailedToRemoveUser struct {
	User  User
	Error error
}
