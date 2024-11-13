package main

import (
	"context"
	"fmt"
	user_domain "go-boilerplate/internal/app/module/user/domain"
	"go-boilerplate/internal/pkg/infrastructure/kernel"
	"go-boilerplate/pkg/config"
	"go-boilerplate/test"
)

func main() {
	// Load the configuration from the .env file
	cnf := config.LoadConfig()

	// Initialize kernel with configuration
	k := kernel.Init(cnf)

	userModule := k.Modules["user_module"].(*kernel.UserModule)

	CreateAndSaveUsers(userModule.UserRepository)
}

func CreateAndSaveUsers(repository user_domain.UserRepository) error {
	for i := 0; i < 500; i++ {
		// Generate a random user using the fixture function
		user, err := test.RandomUserFixture()
		if err != nil {
			return fmt.Errorf("failed to create user fixture: %w", err)
		}

		// Save the user to the repository
		err = repository.Save(context.Background(), user)
		if err != nil {
			return fmt.Errorf("failed to save user %d: %w", i, err)
		}
	}
	return nil
}
