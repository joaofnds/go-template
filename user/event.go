package user

type UserCreated struct {
	User User `json:"user"`
}

func NewUserCreated(u User) UserCreated {
	return UserCreated{User: u}
}
