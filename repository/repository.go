package repository

import (
	"context"
	"errors"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/windbnb/user-service/model"
	"github.com/windbnb/user-service/tracer"
)

type IRepository interface {
	CheckCredentials(email, password string, ctx context.Context) (model.User, error)
	CreateUser(user model.User, ctx context.Context) (model.User, error)
	FindUserById(id uint64, ctx context.Context) (model.User, error)
	SaveUser(user model.User, ctx context.Context) (model.User, error)
	SaveUserDeletionEvent(userId uint64, ctx context.Context)
	DeleteUser(userId uint64, ctx context.Context) error
}

type Repository struct {
	Db *gorm.DB
}

func (r *Repository) CheckCredentials(email string, password string, ctx context.Context) (model.User, error) {
	span := tracer.StartSpanFromContext(ctx, "checkCredentialsRepository")
	defer span.Finish()

	var user model.User

	r.Db.Table("users").Where("email = ? AND password = ?", email, password).First(&user)

	if user.ID == 0 {
		err := errors.New("user does not exist")
		tracer.LogError(span, err)
		return user, err
	}

	return user, nil
}

func (r *Repository) CreateUser(user model.User, ctx context.Context) (model.User, error) {
	span := tracer.StartSpanFromContext(ctx, "createUserRepository")
	defer span.Finish()

	createdUser := r.Db.Create(&user)

	if createdUser.Error != nil {
		tracer.LogError(span, createdUser.Error)
		return user, createdUser.Error
	}

	return user, nil
}

func (r *Repository) FindUserById(id uint64, ctx context.Context) (model.User, error) {
	span := tracer.StartSpanFromContext(ctx, "findUserByIdRepository")
	defer span.Finish()

	var user model.User

	r.Db.First(&user, id)

	if user.ID == 0 {
		err := errors.New("there is no user with id " + strconv.FormatUint(uint64(id), 10))
		tracer.LogError(span, err)
		return model.User{}, err
	}

	return user, nil
}

func (r *Repository) SaveUser(user model.User, ctx context.Context) (model.User, error) {
	span := tracer.StartSpanFromContext(ctx, "saveUserRepository")
	defer span.Finish()

	createdUser := r.Db.Save(&user)

	if createdUser.Error != nil {
		tracer.LogError(span, createdUser.Error)
		return user, createdUser.Error
	}

	return user, nil
}

func (r *Repository) SaveUserDeletionEvent(userId uint64, ctx context.Context) {
	span := tracer.StartSpanFromContext(ctx, "saveUserDeletionEventRepository")
	defer span.Finish()

	r.Db.Save(&model.UserDeletionEvent{UserId: userId})
}

func (r *Repository) DeleteUser(userId uint64, ctx context.Context) error {
	span := tracer.StartSpanFromContext(ctx, "deleteUserRepository")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	userToDelete, err := r.FindUserById(userId, ctx)
	if err != nil {
		tracer.LogError(span, err)
		return err
	}

	r.Db.Delete(userToDelete)

	return nil
}
