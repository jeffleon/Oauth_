package repository

import (
	"github.com/jeffleon/oauth-microservice/pkg/user/domain"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserDBRepository {
	return &userRepository{
		db: db,
	}
}

func (ur *userRepository) GetUser(id int64) (*domain.User, error) {
	var user domain.User
	if err := ur.db.Where("id = ?", id).Find(&user).Error; err != nil {
		logrus.Errorf("FindByID error getting row UserID: %d, error: %s", id, err.Error())
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) CreateUser(user *domain.User) (*domain.User, error) {
	hashedPassword, err := user.Hash()
	if err != nil {
		logrus.Errorf("User create Error: %s", err.Error())
		return nil, err
	}
	user.Password = string(hashedPassword)
	user.Active = true
	if err := ur.db.Create(&user).Error; err != nil {
		logrus.Errorf("User create Error: %s", err.Error())
		return nil, err
	}
	logrus.Infof("User %d created successfully", user.ID)
	return user, nil
}

func (ur *userRepository) UpdateUser(user *domain.User) error {
	_, err := ur.GetUser(user.ID)
	if err != nil {
		logrus.Errorf("User %d update Error: %s", user.ID, err.Error())
		return err
	}
	if user.Password != "" {
		hashedPassword, err := user.Hash()
		if err != nil {
			logrus.Errorf("User %d update Error: %s", user.ID, err.Error())
			return err
		}
		user.Password = string(hashedPassword)
	}
	if err := ur.db.Where("id = ?", user.ID).Updates(&user).Error; err != nil {
		logrus.Errorf("User %d update Error: %s", user.ID, err.Error())
		return err
	}
	logrus.Infof("User %d updated successfully", user.ID)
	return nil
}

func (ur *userRepository) DeleteUser(id int64) (*domain.User, error) {
	var user domain.User
	if err := ur.db.Where("id = ?", id).Delete(&user).Error; err != nil {
		logrus.Errorf("User %d delete Error: %s", id, err.Error())
		return nil, err
	}
	logrus.Infof("User %d deleted successfully", id)
	return &user, nil
}

func (ur *userRepository) VerifyPassword(password, hashedPassword []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

func (ur *userRepository) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := ur.db.Where("email = ?", email).Find(&user).Error; err != nil {
		logrus.Errorf("FindByID error getting row email: %s, error: %s", email, err.Error())
		return nil, err
	}
	return &user, nil
}
