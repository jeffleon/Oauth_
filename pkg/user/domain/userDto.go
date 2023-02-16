package domain

import "strings"

type Userdto struct {
	ID           int64       `json:"id,omitempty"`
	FirstName    string      `json:"first_name"`
	LastName     string      `json:"last_name"`
	Email        string      `json:"email"`
	Password     string      `json:"password,omitempty"`
	Phone        PhoneStruct `json:"phone"`
	Token        string      `json:"token,omitempty"`
	RefreshToken string      `json:"refresh_token,omitempty"`
}

func (u Userdto) UserDto2Domain() *User {
	return &User{
		ID:          u.ID,
		FirstName:   strings.ToLower(u.FirstName),
		LastName:    strings.ToLower(u.LastName),
		Email:       strings.ToLower(u.Email),
		Password:    u.Password,
		Phone:       u.Phone.Number,
		CountryCode: u.Phone.CountryCode,
	}
}

type PhoneStruct struct {
	Number      string `json:"number" bson:"number" structs:"number"`
	CountryCode string `json:"country_code" bson:"country_code" structs:"country_code"`
}

type SignIndto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Tokendto struct {
	ID           int64  `json:"id"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type ForgotPassworddto struct {
	Email string `json:"email"`
}
