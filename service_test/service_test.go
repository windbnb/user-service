package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/windbnb/user-service/model"
	"github.com/windbnb/user-service/repository"
	"github.com/windbnb/user-service/service"
)

func TestLogin_SuccessfulLogin(t *testing.T) {
	// Create a mock repository with desired behavior
	mockRepo := &MockRepo{
		CheckCredentialsFn: func(email, password string, ctx context.Context) (model.User, error) {
			// Return a mock user for successful login
			return model.User{
				Email: "test@example.com",
				Password: "password",
				Role:  model.HOST,
			}, nil
		},
	}

	// Create an instance of the UserService with the mock repository
	userService := service.UserService{
		Repo: mockRepo,
	}

	// Call the Login function with mock credentials and context
	credentials := model.Credentials{
		Email:    "test@example.com",
		Password: "password",
	}
	token, err := userService.Login(credentials, context.Background())

	// Assert that the returned token is not empty and there is no error
	assert.NotEmpty(t, token)
	assert.NoError(t, err)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	// Create a mock repository with desired behavior
	mockRepo := &MockRepo{
		CheckCredentialsFn: func(email, password string, ctx context.Context) (model.User, error) {
			// Return an error for invalid credentials
			return model.User{}, errors.New("invalid credentials")
		},
	}

	// Create an instance of the UserService with the mock repository
	userService := service.UserService{
		Repo: mockRepo,
	}

	// Call the Login function with mock credentials and context
	credentials := model.Credentials{
		Email:    "test@example.com",
		Password: "password",
	}
	token, err := userService.Login(credentials, context.Background())

	// Assert that the returned token is empty and there is an error
	assert.Empty(t, token)
	assert.EqualError(t, err, "bad credentials")
}


type MockRepo struct {
	repository.Repository
	CheckCredentialsFn func(email, password string, ctx context.Context) (model.User, error)
}

func (m *MockRepo) CheckCredentials(email, password string, ctx context.Context) (model.User, error) {
	return m.CheckCredentialsFn(email, password, ctx)
}

func (m *MockRepo) CreateUser(user model.User, ctx context.Context) (model.User, error) {
	return user, nil
}

func (m *MockRepo) FindUserById(id uint64, ctx context.Context) (model.User, error) {
	return model.User{}, nil
}

func (m *MockRepo) SaveUser(user model.User, ctx context.Context) (model.User, error) {
	return user, nil
}

func (m *MockRepo) SaveUserDeletionEvent(userId uint64, ctx context.Context) {
}

func (m *MockRepo) DeleteUser(userId uint64, ctx context.Context) error {
	return nil
}
