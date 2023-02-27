package aplication

import (
	"time"

	"github.com/jeffleon/oauth-microservice/internal/config"
	"github.com/jeffleon/oauth-microservice/pkg/user/domain"
	"github.com/sirupsen/logrus"
)

type UserService interface {
	GetUser(int64) (*domain.Userdto, error)
	SingUp(*domain.Userdto) (*domain.Userdto, error)
	UpdateUser(*domain.Userdto, string) error
	DeleteUser(int64) (*domain.Userdto, error)
	Logout(string) error
	SignIn(*domain.SignIndto) (*domain.Tokendto, error)
	RefreshAccessToken(int64) (*domain.Tokendto, error)
	VerifyToken(string) (*domain.Claims, bool)
	VerifyBlackList(token string) bool
}

type userService struct {
	userRepository  domain.UserDBRepository
	tokenRepository domain.TokenRepository
	redisRepository domain.RedisRepository
	rpcRepository   domain.RPCRepository
}

func NewUserService(userRepository domain.UserDBRepository,
	tokenRepository domain.TokenRepository,
	redisRepository domain.RedisRepository,
	rpcRepository domain.RPCRepository,
) UserService {
	return &userService{
		userRepository,
		tokenRepository,
		redisRepository,
		rpcRepository,
	}
}

func (us *userService) GetUser(id int64) (*domain.Userdto, error) {
	user, err := us.userRepository.GetUser(id)
	if err != nil {
		return nil, err
	}
	userDto := user.UserDomain2Dto()
	return userDto, nil
}

func (us *userService) RefreshAccessToken(id int64) (*domain.Tokendto, error) {
	user, err := us.userRepository.GetUser(id)
	if err != nil {
		return nil, err
	}

	token, err := us.tokenRepository.CreateToken(user, "predeterminated")
	if err != nil {
		return nil, err
	}

	return &domain.Tokendto{Token: *token, ID: user.ID}, nil
}

func (us *userService) SingUp(user *domain.Userdto) (*domain.Userdto, error) {
	userDomain := user.UserDto2Domain()
	response, err := us.userRepository.CreateUser(userDomain)
	if err != nil {
		return nil, err
	}

	token, err := us.tokenRepository.CreateToken(response, "predeterminated")
	if err != nil {
		return nil, err
	}

	refreshToken, err := us.tokenRepository.CreateToken(response, "refreshToken")
	if err != nil {
		return nil, err
	}
	user.ID = response.ID
	user.RefreshToken = *refreshToken
	user.Token = *token

	payload := domain.NewRPCPayloadEmailWelcome(user.Email)

	res, err := us.rpcRepository.SendEmail(payload)
	if err != nil {
		logrus.Errorf("Error sent email %s", err)
	} else {
		logrus.Infof("Email OK, %s", res)
	}

	return user, nil
}

func (us *userService) UpdateUser(user *domain.Userdto, updateType string) error {
	userDomain := user.UserDto2Domain()
	if updateType == "user" {
		userDomain.Password = ""
	}
	return us.userRepository.UpdateUser(userDomain)
}

func (us *userService) DeleteUser(id int64) (*domain.Userdto, error) {
	user, err := us.userRepository.DeleteUser(id)
	if err != nil {
		return nil, err
	}
	userDto := user.UserDomain2Dto()
	return userDto, nil
}

func (us *userService) SignIn(user *domain.SignIndto) (*domain.Tokendto, error) {
	response, err := us.userRepository.GetUserByEmail(user.Email)
	if err != nil {
		return nil, err
	}

	err = us.userRepository.VerifyPassword([]byte(user.Password), []byte(response.Password))
	if err != nil {
		return nil, err
	}

	token, err := us.tokenRepository.CreateToken(response, "predeterminated")
	if err != nil {
		return nil, err
	}

	return &domain.Tokendto{ID: response.ID, Token: *token}, nil
}

func (us *userService) Logout(token string) error {
	return us.redisRepository.HSet(config.Config.TokenBlackListSet, token, time.Now().Format("2017-09-07 17:06:06"))
}

func (us *userService) VerifyBlackList(token string) bool {
	valid, err := us.redisRepository.HGet(config.Config.TokenBlackListSet, token)
	if (err != nil && err.Error() != "redis: nil") || valid != "" {
		logrus.Errorf("Invalid token since %s or with error %s", valid, err)
		return false
	}

	return true
}

func (us *userService) VerifyToken(token string) (*domain.Claims, bool) {
	claims, err := us.tokenRepository.VerifyToken(token, "predetermined")

	return claims, err == nil
}
