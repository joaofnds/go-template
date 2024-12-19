package user_http

type UserCreateDTO struct {
	Email string `json:"email" validate:"required,email"`
}
