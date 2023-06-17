package model

import (
	"github.com/jinzhu/gorm"
)

type UserRole string

const (
	HOST  UserRole = "HOST"
	GUEST UserRole = "GUEST"
)

type User struct {
	gorm.Model
	Email string `gorm:"not null;default:null;unique_index"`
	Username string `gorm:"not null;default:null;unique_index"`
	Password string `gorm:"not null;default:null"`
	Name string `gorm:"not null;default:null"`
	Surname string `gorm:"not null;default:null"`
	Address string `gorm:"not null;default:null"`
	Role UserRole `gorm:"not null;default:null"`
	ReservationRequestNotification bool `gorm:"not null;default:false"`
	ReservationCanceledNotification bool `gorm:"not null;default:false"`
	SelfReviewNotification bool `gorm:"not null;default:false"`
	AccomodationReviewNotification bool `gorm:"not null;default:false"`
	ReservationStatusChangedNotification bool `gorm:"not null;default:false"`
}

func (user *User) ToDTO() UserResponseDTO {
	return UserResponseDTO{Id: user.ID, Email: user.Email, Name: user.Name, Surname: user.Surname, Address: user.Address, Username: user.Username, Role: user.Role,
							ReservationRequestNotification: user.ReservationRequestNotification, 
							ReservationCanceledNotification: user.ReservationCanceledNotification, 
							SelfReviewNotification: user.SelfReviewNotification,
							AccomodationReviewNotification: user.AccomodationReviewNotification, 
							ReservationStatusChangedNotification: user.ReservationStatusChangedNotification}
}

type UserDeletionEvent struct {
	gorm.Model
	UserId uint64 `gorm:"not null;default:null"`
}
