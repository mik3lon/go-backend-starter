package user_domain

type UserAlreadyExists struct {
	extraItems map[string]interface{}
}

func NewUserAlreadyExists(email string) *UserAlreadyExists {
	return &UserAlreadyExists{
		extraItems: map[string]interface{}{
			"email": email,
		},
	}
}

func (u UserAlreadyExists) Error() string {
	return "user already exists"
}

func (u UserAlreadyExists) ExtraItems() map[string]interface{} {
	return map[string]interface{}{}
}
