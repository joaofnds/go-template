package user

type UserService struct {
	users []User
}

func NewUserService() *UserService {
	return &UserService{}
}

func (service *UserService) CreateUser(name string) User {
	user := User{name}
	service.users = append(service.users, user)
	return user
}

func (service *UserService) DeleteAll() {
	service.users = service.users[:0]
}

func (service *UserService) List() []User {
	return service.users
}

func (service *UserService) FindByName(name string) (User, bool) {
	for _, u := range service.users {
		if u.Name == name {
			return u, true
		}
	}

	return User{}, false
}

func (service *UserService) Remove(user User) bool {
	index := -1

	for i, u := range service.users {
		if u == user {
			index = i
			break
		}
	}

	if index == -1 {
		return false
	}

	service.users = append(service.users[:index], service.users[index+1:]...)

	return true
}
