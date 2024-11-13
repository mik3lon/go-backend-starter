package user_domain

import (
	"time"
)

type UserList []*User

type User struct {
	ID                string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"` // UUID primary key
	Username          string    `gorm:"type:varchar(50);uniqueIndex"`                    // Unique index for username
	Email             string    `gorm:"type:varchar(100);uniqueIndex"`                   // Unique index for email
	HashedPassword    string    `gorm:"type:varchar(255)"`                               // Password hash
	Name              string    `gorm:"type:varchar(50)"`                                // Name field
	Surname           string    `gorm:"type:varchar(50)"`                                // Surname field
	Role              string    `gorm:"type:varchar(20);default:'user'"`                 // Role with a default value
	ProfilePictureUrl string    `gorm:"type:varchar(200)"`                               // profile picture url
	CreatedAt         time.Time `gorm:"autoCreateTime"`                                  // Automatically set on create
	UpdatedAt         time.Time `gorm:"autoUpdateTime"`                                  // Automatically set on update
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
