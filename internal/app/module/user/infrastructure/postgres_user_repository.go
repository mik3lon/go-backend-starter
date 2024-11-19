package user_infrastructure

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	user_domain "github.com/mik3lon/starter-template/internal/app/module/user/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgresUserRepository is a Postgres implementation of UserRepository using Gorm.
type PostgresUserRepository struct {
	DB *gorm.DB
}

// NewPostgresUserRepository initializes a new Postgres user repository.
func NewPostgresUserRepository(dsn string) (*PostgresUserRepository, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DriverName: "pgx",
		DSN:        dsn,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Ensure the User table exists
	if err = db.AutoMigrate(&user_domain.User{}); err != nil {
		return nil, err
	}

	return &PostgresUserRepository{
		DB: db,
	}, nil
}

func (r *PostgresUserRepository) Save(ctx context.Context, user *user_domain.User) error {
	result := r.DB.WithContext(ctx).Save(user)
	if result.Error != nil {
		// Check if the error is a unique constraint violation
		var pgErr *pgconn.PgError
		errors.As(result.Error, &pgErr)
		if pgErr.Code == "23505" {
			return user_domain.NewUserAlreadyExists(user.Email)
		}
		return fmt.Errorf("failed to save user: %w", result.Error)
	}
	return nil
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*user_domain.User, error) {
	var user *user_domain.User
	result := r.DB.First(&user, "email = ?", email)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, user_domain.NewUserNotFound(email)
	}
	return user, result.Error
}

func (r *PostgresUserRepository) FindAll(ctx context.Context, page int, size int) (user_domain.UserList, error) {
	var users []*user_domain.User
	offset := (page - 1) * size

	result := r.DB.Limit(size).Offset(offset).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	userList := make(user_domain.UserList, len(users))
	for i, u := range users {
		userList[i] = u
	}
	return userList, nil
}
