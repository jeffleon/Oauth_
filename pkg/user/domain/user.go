package domain

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserDBRepository interface {
	GetUser(int64) (*User, error)
	CreateUser(*User) (*User, error)
	UpdateUser(*User) error
	DeleteUser(int64) (*User, error)
	GetUserByEmail(string) (*User, error)
	VerifyPassword(password, hashedPassword []byte) error
}

type User struct {
	gorm.Model
	ID          int64  `gorm:"primary_key;auto_increment" json:"id"`
	FirstName   string `gorm:"fist_name"`
	LastName    string `gorm:"last_name"`
	Email       string `gorm:"email;size:255;not null;unique"`
	Password    string `gorm:"password;not null"`
	Phone       string `gorm:"phone;size:255;not null"`
	CountryCode string `gorm:"country_code"`
	Active      bool   `gorm:"active"`
}

func (u *User) UserDomain2Dto() *Userdto {
	return &Userdto{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Phone: PhoneStruct{
			Number:      u.Phone,
			CountryCode: u.CountryCode,
		},
	}
}

func (u *User) Hash() ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
}

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&User{}); err != nil {
		return err
	}
	return nil
}
