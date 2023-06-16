package model

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type UserDTO struct {
	Id       uint   `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Address  string `json:"address"`
	Username string `json:"username"`
}

func (user *UserDTO) ToUser() User {
	return User{Email: user.Email, Name: user.Name, Surname: user.Surname, Address: user.Address}
}

type UserResponseDTO struct {
	Id       uint     `json:"id"`
	Email    string   `json:"email"`
	Name     string   `json:"name"`
	Surname  string   `json:"surname"`
	Address  string   `json:"address"`
	Username string   `json:"username"`
	Role     UserRole `json:"role"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	Email string   `json:"email"`
	Role  UserRole `json:"role"`
	Id    uint     `json:"id"`
	jwt.StandardClaims
}

type ErrorResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type CreateUserRequest struct {
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Name     string   `json:"name"`
	Surname  string   `json:"surname"`
	Address  string   `json:"address"`
	Role     UserRole `json:"role"`
}

func (user *CreateUserRequest) ToUser() User {
	return User{Email: user.Email, Name: user.Name, Surname: user.Surname, Password: user.Password, Address: user.Address, Role: user.Role, Username: user.Username}
}

type ReservationRequestStatus string

const (
	SUBMITTED ReservationRequestStatus = "SUBMITTED"
	ACCEPTED  ReservationRequestStatus = "ACCEPTED"
	DECLINED  ReservationRequestStatus = "DECLINED"
)

type ReservationRequestDto struct {
	Status               ReservationRequestStatus `json:"status"`
	GuestID              uint                     `json:"guestID"`
	AccommodationID      uint                     `json:"accommodationID"`
	ReservationRequestID uint                     `json:"reservationRequestID"`
	StartDate            time.Time                `json:"startDate"`
	EndDate              time.Time                `json:"endDate"`
	GuestNumber          uint                     `json:"guestNumber"`
}

type ChangePasswordDTO struct {
	Id          uint   `json:"id"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}
