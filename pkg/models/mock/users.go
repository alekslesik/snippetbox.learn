package mock

import (
	"time"

	"github.com/alekslesik/snippetbox.learn/pkg/models"
)

var mockUser = &models.User {
	ID:             1,
	Name:           "Alex",
	Email:          "alekslesik@gmail.com",
	HashedPassword: []byte("password"),
	Created:        time.Now(),
}

type UserModel struct{}

// Rewrite all mysql.UserModel methods

// Add a new record to the users table.
func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

// Verify whether a user exists with the provided email address and password.
// Return the relevant user ID if they do.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	switch {
	case email == mockUser.Email && password == string(mockUser.HashedPassword):
		return mockUser.ID, nil
	default:
		return 0, models.ErrInvalidCredentials
	}
}

// Fetch details for a specific user based on their user ID.
func (m *UserModel) Get(id int) (*models.User, error) {
	switch id {
	case 1:
		return mockUser, nil
	default:
		return nil, models.ErrNoRecord
	}
}
