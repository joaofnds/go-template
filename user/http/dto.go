package http

type UserCreateDTO struct {
	Name string `json:"name" validate:"required,min=3"`
}
