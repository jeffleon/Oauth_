package infraestructure

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jeffleon/oauth-microservice/pkg/user/domain"
	"github.com/sirupsen/logrus"
)

type TokenRepository struct {
	tokenObj domain.TokenObj
}

func NewTokenRepository(tokenObj domain.TokenObj) domain.TokenRepository {
	return &TokenRepository{
		tokenObj: tokenObj,
	}
}

func (t TokenRepository) ChooseSecretExp(tokenType string) ([]byte, time.Duration) {
	switch tokenType {
	case "predetermined":
		return []byte(t.tokenObj.GetPredeterminedSecret()), t.tokenObj.GetPredeterminedExpiration()
	case "refreshToken":
		return []byte(t.tokenObj.GetRefreshTokenSecret()), t.tokenObj.GetRefreshTokenExpiration()
	default:
		return []byte(t.tokenObj.GetPredeterminedSecret()), t.tokenObj.GetPredeterminedExpiration()
	}
}

func (t TokenRepository) CreateToken(user *domain.User, tokenType string) (*string, error) {
	secretKey, expiration := t.ChooseSecretExp(tokenType)
	claims := domain.Claims{
		Email:      user.Email,
		UserID:     user.ID,
		Authorized: true,
		Exp:        time.Now().Add(expiration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		logrus.Errorf("Create token Error: %s", err.Error())
		return nil, err
	}
	logrus.Infof("Token created")
	return &signedToken, nil
}

func (t TokenRepository) VerifyToken(receivedToken string, tokenType string) (*domain.Claims, error) {
	secretKey, _ := t.ChooseSecretExp(tokenType)
	token, err := jwt.Parse(receivedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return secretKey, nil
	})

	if token == nil {
		return nil, fmt.Errorf("non-decryptable token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &domain.Claims{
			Email:  claims["email"].(string),
			UserID: int64(claims["user_id"].(float64)), // TODO TASK: test bigints numbers parsing float to int
		}, nil
	}
	return nil, err
}

func (t TokenRepository) GetToken() {}
