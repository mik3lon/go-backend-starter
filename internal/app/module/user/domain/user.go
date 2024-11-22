package user_domain

import (
	"time"
)

type UserList []*User

type User struct {
	ID                string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Username          string    `gorm:"type:varchar(50);uniqueIndex"`
	Email             string    `gorm:"type:varchar(100);uniqueIndex"`
	HashedPassword    string    `gorm:"type:varchar(255)"`
	Name              string    `gorm:"type:varchar(50)"`
	Surname           string    `gorm:"type:varchar(50)"`
	Role              string    `gorm:"type:varchar(20);default:'user'"`
	ProfilePictureUrl string    `gorm:"type:varchar(200)"`
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime"`
}

func (u *User) UpdateProfile(username string, name string, surname string) {
	u.Username = username
	u.Name = name
	u.Surname = surname
}

func (u *User) UpdateProfilePhoto(image string) {
	u.ProfilePictureUrl = image
}

// CreateUser creates a new User entity.
func CreateUser(id, username, email, password, name, surname, role, profilePictureUrl string) *User {
	return &User{
		ID:                id,
		Username:          username,
		Email:             email,
		HashedPassword:    password,
		Name:              name,
		Surname:           surname,
		Role:              role,
		ProfilePictureUrl: profilePictureUrl,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
}

func FromPrimitives(
	id,
	username,
	email,
	hashedPassword,
	name,
	surname,
	role,
	profilePictureUrl string,
	createdAt, updatedAt time.Time,
) *User {
	return &User{
		ID:                id,
		Username:          username,
		Email:             email,
		HashedPassword:    hashedPassword,
		Name:              name,
		Surname:           surname,
		Role:              role,
		ProfilePictureUrl: profilePictureUrl,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}
}
