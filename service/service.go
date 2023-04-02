package service

import (
	"errors"
	"net/mail"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/windbnb/user-service/model"
	"github.com/windbnb/user-service/repository"
)

var jwtKey = []byte("z7031Q8Qy9zVO-T2o7lsFIZSrd05hH0PaeaWIBvLh9s")

type UserService struct {
	Repo *repository.Repository
}

func (service *UserService) Login(credentials model.Credentials) (string, error) {
	user, err := service.Repo.CheckCredentials(credentials.Email, credentials.Password)

	if err != nil {
		return "", errors.New("bad credentials")
	}

	expirationTime := time.Now().Add(time.Hour * 24)
	claims := model.Claims{Email: user.Email, Role: user.Role, Id: user.ID, StandardClaims: jwt.StandardClaims{ExpiresAt: expirationTime.Unix()}}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	tokenString, _ := token.SignedString(jwtKey)

	return tokenString, nil
}

func (service *UserService) CreateUser(user model.User) (model.User, error) {
	_, err := mail.ParseAddress(user.Email)

	if err != nil {
		return user, errors.New("email format is not valid")
	}

	createdUser, err := service.Repo.CreateUser(user)

	if err != nil {
		return user, errors.New("error while trying to save user")
	}

	return createdUser, nil
}

func (service *UserService) AuthoriseUser(tokenString string, role model.UserRole) error {

	claims := model.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, 
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil || !token.Valid {
		return err
	}

	if token.Claims.(*model.Claims).Role != role {
		return errors.New("user does not have said role")
	}

	return nil
}

func (service *UserService) FindUser(userId uint64) (model.User, error) {
	userToUpdate, err := service.Repo.FindUserById(userId)

	if err != nil {
		return model.User{}, errors.New("user with given id does not exist")
	}

	return userToUpdate, nil
}

func (service *UserService) EditUser(user model.UserDTO, userId uint64) error {
	userToUpdate, err := service.Repo.FindUserById(userId)

	if err != nil {
		return errors.New("user with given id does not exist")
	}

	userToUpdate.Name = user.Name
	userToUpdate.Surname = user.Surname
	userToUpdate.Email = user.Email
	userToUpdate.Address = user.Address

	if user.OldPassword != "" {
		if user.OldPassword != userToUpdate.Password {
			return errors.New("old and new password do not match")
		}

		userToUpdate.Password = user.NewPassword
	}

	_, err = service.Repo.SaveUser(userToUpdate)

	if err != nil {
		return errors.New("error while saving user")
	}

	return nil
}
