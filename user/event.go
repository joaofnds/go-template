package user

type UserCreated struct {
	User User
}

type UserRemoved struct {
	User User
}

type FailedToCreateUser struct {
	Error string
}

type FailedToDeleteAll struct {
	Error string
}

type FailedToFindByID struct {
	Error string
}

type FailedToFindByName struct {
	Error string
}

type FailedToRemoveUser struct {
	User  User
	Error string
}
