package test

import (
	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
	user_domain "go-boilerplate/internal/app/module/user/domain"
	"math/rand"
	"time"
)

const dateFormat = "2006-01-02" // Format used by faker.Date()

// RandomUserFixture creates a randomized User using faker and the FromPrimitives method.
func RandomUserFixture() (*user_domain.User, error) {
	// Generate random data
	id := uuid.New().String()          // Generate a unique ID
	username := faker.Username()       // Generate a random username
	email := faker.Email()             // Generate a random email
	hashedPassword := faker.Password() // Generate a random password (hashed)
	name := faker.FirstName()          // Generate a random first name
	surname := faker.LastName()        // Generate a random last name
	role := randomRole()               // Randomly select a role
	profilePictureUrl := faker.URL()
	createdAtStr := faker.Date()
	createdAt, err := time.Parse(dateFormat, createdAtStr)
	if err != nil {
		return nil, err
	}

	updatedAt := createdAt.Add(24 * time.Hour) // Add 24 hours to createdAt for updatedAt

	// Create the User using FromPrimitives
	user := user_domain.FromPrimitives(
		id,
		username,
		email,
		hashedPassword,
		name,
		surname,
		role,
		profilePictureUrl,
		createdAt,
		updatedAt,
	)

	return user, nil
}

// randomRole selects a random role from a predefined list
func randomRole() string {
	roles := []string{"user", "admin", "moderator"}
	rand.Seed(time.Now().UnixNano()) // Seed to ensure randomness
	return roles[rand.Intn(len(roles))]
}
