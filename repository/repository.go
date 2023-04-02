package repository

import (
	"errors"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/windbnb/user-service/model"
)

type Repository struct {
	Db *gorm.DB
}

func (r *Repository) CheckCredentials(email string, password string) (model.User, error) {
	var user model.User

	r.Db.Table("users").Where("email = ? AND password = ?", email, password).First(&user)

	if user.ID == 0 {
		return user, errors.New("user does not exist")
	}

	return user, nil
}

func (r *Repository) CreateUser(user model.User) (model.User, error) {
	createdUser := r.Db.Create(&user)

	if createdUser.Error != nil {
		return user, createdUser.Error
	}

	return user, nil
}

func (r *Repository) FindUserById(id uint64) (model.User, error) {
	var user model.User

	r.Db.First(&user, id)

	if user.ID == 0 {
		return model.User{}, errors.New("there is no user with id " + strconv.FormatUint(uint64(id), 10))
	}

	return user, nil
}

func (r *Repository) SaveUser(user model.User) (model.User, error) {
	createdUser := r.Db.Save(&user)

	if createdUser.Error != nil {
		return user, createdUser.Error
	}

	return user, nil
}

func (r *Repository) DeleteUser() {

}