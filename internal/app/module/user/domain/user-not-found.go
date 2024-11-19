package user_domain

type UserNotFound struct {
	extraItems map[string]interface{}
}

func NewUserNotFound(email string) *UserNotFound {
	return &UserNotFound{
		extraItems: map[string]interface{}{
			"email": email,
		},
	}
}

func (u UserNotFound) Error() string {
	return "user not found"
}

func (u UserNotFound) ExtraItems() map[string]interface{} {
	return map[string]interface{}{}
}
