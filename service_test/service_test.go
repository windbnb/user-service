package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/windbnb/user-service/model"
	"github.com/windbnb/user-service/repository"
	"github.com/windbnb/user-service/service"
	"github.com/windbnb/user-service/util"
)

func TestLogin_ValidCredentials_Integration(t *testing.T) {
	db := util.ConnectToDatabase()
	defer db.Close()
	userService := service.UserService{Repo: &repository.Repository{Db: db}}

	credentials := model.Credentials {
		Email:    "host@email.com",
		Password: "host",
	}

	token, err := userService.Login(credentials, context.Background())

	assert.NotEmpty(t, token)
	assert.NoError(t, err)
}

func TestLogin_InvalidCredentials_Integration(t *testing.T) {
	db := util.ConnectToDatabase()
	defer db.Close()
	userService := service.UserService{Repo: &repository.Repository{Db: db}}

	credentials := model.Credentials{
		Email:    "host@email.com",
		Password: "wrong_password",
	}

	token, err := userService.Login(credentials, context.Background())

	assert.Empty(t, token)
	assert.EqualError(t, err, "bad credentials")
}

func TestLogin_SuccessfulLogin(t *testing.T) {
	mockRepo := &MockRepo{
		CheckCredentialsFn: func(email, password string, ctx context.Context) (model.User, error) {
			return model.User{
				Email: "test@example.com",
				Password: "password",
				Role:  model.HOST,
			}, nil
		},
	}

	userService := service.UserService{
		Repo: mockRepo,
	}

	credentials := model.Credentials{
		Email:    "test@example.com",
		Password: "password",
	}
	token, err := userService.Login(credentials, context.Background())

	assert.NotEmpty(t, token)
	assert.NoError(t, err)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	mockRepo := &MockRepo{
		CheckCredentialsFn: func(email, password string, ctx context.Context) (model.User, error) {
			return model.User{}, errors.New("invalid credentials")
		},
	}

	userService := service.UserService{
		Repo: mockRepo,
	}

	credentials := model.Credentials{
		Email:    "test@example.com",
		Password: "password",
	}
	token, err := userService.Login(credentials, context.Background())

	assert.Empty(t, token)
	assert.EqualError(t, err, "bad credentials")
}

func TestCreateUser_InvalidEmailFormat(t *testing.T) {
	mockRepo := &MockRepo{}

	userService := service.UserService{
		Repo: mockRepo,
	}

	user := model.User{
		Email: "invalid_email",
		ReservationRequestNotification:true, 
		ReservationCanceledNotification:true, 
		SelfReviewNotification:true, 
		AccomodationReviewNotification:true,
	}

	createdUser, err := userService.CreateUser(user, context.Background())

	assert.Equal(t, user, createdUser)
	assert.EqualError(t, err, "email format is not valid")
}

func TestCreateUser_ErrorSavingUser(t *testing.T) {
	mockRepo := &MockRepo{
		CreateUserFn: func(user model.User, ctx context.Context) (model.User, error) {
			return model.User{}, errors.New("error saving user")
		},
	}

	userService := service.UserService{
		Repo: mockRepo,
	}

	user := model.User{
		Email: "test@example.com",
		ReservationRequestNotification:true, 
		ReservationCanceledNotification:true, 
		SelfReviewNotification:true, 
		AccomodationReviewNotification:true,
	}

	createdUser, err := userService.CreateUser(user, context.Background())

	assert.Equal(t, user, createdUser)
	assert.EqualError(t, err, "error while trying to save user")
}

func TestCreateUser_Successful(t *testing.T) {
	mockRepo := &MockRepo{
		CreateUserFn: func(user model.User, ctx context.Context) (model.User, error) {
			return user, nil
		},
	}

	userService := service.UserService{
		Repo: mockRepo,
	}

	user := model.User{
		Email: "test@example.com",
		ReservationRequestNotification:true, 
		ReservationCanceledNotification:true, 
		SelfReviewNotification:true, 
		AccomodationReviewNotification:true,
	}

	createdUser, err := userService.CreateUser(user, context.Background())

	assert.Equal(t, user, createdUser)
	assert.NoError(t, err)
}


type MockRepo struct {
	repository.Repository
	CheckCredentialsFn func(email, password string, ctx context.Context) (model.User, error)
	CreateUserFn func(user model.User, ctx context.Context) (model.User, error)
}

func (m *MockRepo) CheckCredentials(email, password string, ctx context.Context) (model.User, error) {
	return m.CheckCredentialsFn(email, password, ctx)
}

func (m *MockRepo) CreateUser(user model.User, ctx context.Context) (model.User, error) {
	return m.CreateUserFn(user, ctx)
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
