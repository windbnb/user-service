package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/mail"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/windbnb/user-service/client"
	"github.com/windbnb/user-service/model"
	"github.com/windbnb/user-service/repository"
	"github.com/windbnb/user-service/tracer"
)


var jwtKey []byte

func InitJWTKey() {
	keyLength := 32

	randomString := make([]byte, keyLength)
	_, err := rand.Read(randomString)
	if err != nil {
		panic(err)
	}

	jwtKey = []byte(base64.RawURLEncoding.EncodeToString(randomString))
}

type UserService struct {
	Repo repository.IRepository
}

func (service *UserService) Login(credentials model.Credentials, ctx context.Context) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "loginService")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	user, err := service.Repo.CheckCredentials(credentials.Email, credentials.Password, ctx)

	if err != nil {
		tracer.LogError(span, err)
		return "", errors.New("bad credentials")
	}

	expirationTime := time.Now().Add(time.Hour * 24)
	claims := model.Claims{Email: user.Email, Role: user.Role, Id: user.ID, StandardClaims: jwt.StandardClaims{ExpiresAt: expirationTime.Unix()}}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	tokenString, _ := token.SignedString(jwtKey)

	return tokenString, nil
}

func (service *UserService) CreateUser(user model.User, ctx context.Context) (model.User, error) {
	span := tracer.StartSpanFromContext(ctx, "createUserService")
	defer span.Finish()

	_, err := mail.ParseAddress(user.Email)

	if err != nil {
		tracer.LogError(span, err)
		return user, errors.New("email format is not valid")
	}

	if user.Role == model.GUEST {
		user.ReservationStatusChangedNotification = true;
	} else {
		user.SelfReviewNotification = true;
		user.AccomodationReviewNotification = true;
		user.ReservationRequestNotification = true;
		user.ReservationCanceledNotification = true;
	}

	ctx = tracer.ContextWithSpan(context.Background(), span)
	createdUser, err := service.Repo.CreateUser(user, ctx)

	if err != nil {
		tracer.LogError(span, err)
		return user, errors.New("error while trying to save user")
	}

	return createdUser, nil
}

func (service *UserService) AuthenticateUser(tokenString string, role model.UserRole, authorise bool, ctx context.Context) (model.User, error) {
	span := tracer.StartSpanFromContext(ctx, "authoriseUserService")
	defer span.Finish()

	claims := model.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, 
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil || !token.Valid {
		tracer.LogError(span, err)
		return model.User{}, err
	}

	if authorise && token.Claims.(*model.Claims).Role != role {
		err := errors.New("user does not have said role")
		tracer.LogError(span, err)
		return model.User{}, err
	}

	ctx = tracer.ContextWithSpan(context.Background(), span)
	user, err := service.Repo.FindUserById(uint64(token.Claims.(*model.Claims).Id), ctx)

	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (service *UserService) FindUser(userId uint64, ctx context.Context) (model.User, error) {
	span := tracer.StartSpanFromContext(ctx, "findUserService")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	userToUpdate, err := service.Repo.FindUserById(userId, ctx)

	if err != nil {
		tracer.LogError(span, err)
		return model.User{}, errors.New("user with given id does not exist")
	}

	return userToUpdate, nil
}

func (service *UserService) DeleteUser(userId uint64, ctx context.Context) error {
	span := tracer.StartSpanFromContext(ctx, "deleteUserService")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	userToDelete, err := service.Repo.FindUserById(userId, ctx)
	if err != nil {
		tracer.LogError(span, err)
		return err
	}

	userIsHost := false
	if userToDelete.Role == model.GUEST {
		err := client.CheckReservations(userToDelete.ID, "guest")
		if err != nil {
			tracer.LogError(span, err)
			return err
		}
	} else if userToDelete.Role == model.HOST {
		userIsHost = true
		err := client.CheckReservations(userToDelete.ID, "owner")
		if err != nil {
			tracer.LogError(span, err)
			return err
		}
	}

	err = service.Repo.DeleteUser(userId, ctx)
	if err != nil {
		tracer.LogError(span, err)
		return err
	}

	if userIsHost {
		err = client.DeleteAccomodationForHost(uint(userId))
		if err != nil {
			tracer.LogError(span, err)
			service.Repo.SaveUserDeletionEvent(userId, ctx)
		}
	}

	return nil
}

func (service *UserService) EditUser(user model.UserDTO, userId uint64, ctx context.Context) error {
	span := tracer.StartSpanFromContext(ctx, "editUserService")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	userToUpdate, err := service.Repo.FindUserById(userId, ctx)

	if err != nil {
		tracer.LogError(span, err)
		return errors.New("user with given id does not exist")
	}

	userToUpdate.Name = user.Name
	userToUpdate.Surname = user.Surname
	userToUpdate.Email = user.Email
	userToUpdate.Address = user.Address

	if userToUpdate.Role == model.GUEST {
		userToUpdate.ReservationStatusChangedNotification = user.ReservationStatusChangedNotification;
	} else {
		userToUpdate.SelfReviewNotification = user.SelfReviewNotification;
		userToUpdate.AccomodationReviewNotification = user.AccomodationReviewNotification;
		userToUpdate.ReservationRequestNotification = user.ReservationRequestNotification;
		userToUpdate.ReservationCanceledNotification = user.ReservationCanceledNotification;
	}

	if user.OldPassword != "" {
		if user.OldPassword != userToUpdate.Password {
			err := errors.New("old and new password do not match")
			tracer.LogError(span, err)
			return err
		}

		userToUpdate.Password = user.NewPassword
	}

	ctx = tracer.ContextWithSpan(context.Background(), span)
	_, err = service.Repo.SaveUser(userToUpdate, ctx)

	if err != nil {
		tracer.LogError(span, err)
		return errors.New("error while saving user")
	}

	return nil
}
