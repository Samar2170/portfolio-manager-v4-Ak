package internal

import (
	"crypto/sha256"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/samar2170/portfolio-manager-v4/internal/models"
)

type RequestUser struct {
	UserID   uint
	Username string
	UserCID  string
}
type UserClaim struct {
	Username string `json:"username"`
	UserCid  string `json:"user_cid"`
	UserID   uint   `json:"user_id"`
	jwt.RegisteredClaims
}

func getCIDForUser() string {
	return uuid.New().String()
}

func generateApiKey() string {
	return uuid.New().String() + "-" + uuid.New().String()
}

func customHash(str string) string {
	hasher := sha256.New()
	hasher.Write([]byte(str))
	bs := hasher.Sum(nil)
	return string(bs)
}

func createToken(u RequestUser) (string, error) {
	claims := UserClaim{
		u.Username,
		u.UserCID,
		u.UserID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "portfolio-manager-ak",
			Subject:   "user token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
func parseToken(token string) (models.User, error) {
	u := models.User{}
	claims := UserClaim{}
	tkn, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		return u, err
	}
	if !tkn.Valid {
		return u, errors.New("invalid token")
	}
	u, err = models.GetUserByCID(claims.UserCid)
	if err != nil {
		return u, err
	}
	return u, nil
}

func GenerateApiKey(token string) (string, error) {
	apiKey := generateApiKey()
	log.Println(apiKey)
	hashedApiKey := customHash(apiKey)
	user, err := parseToken(token)
	if err != nil {
		return "", err
	}
	err = user.UpdateApiKey(hashedApiKey)
	if err != nil {
		return "", err
	}
	return apiKey, nil
}

// func VerifyApiKey(apiKey string) bool {

// }

func Signup(s SignupRequest) error {
	var err error
	dbUser := models.User{
		Username: s.Username,
		UserCID:  getCIDForUser(),
		Password: customHash(s.Password),
		Email:    s.Email,
	}
	err = models.CreateModelInstance(&dbUser)
	if err != nil {
		return err
	}
	err = createGeneralAccountForUser(&dbUser)
	return err
}

func createGeneralAccountForUser(user *models.User) error {
	generalAccount := models.GeneralAccount{
		UserCID:     user.UserCID,
		AccountCode: user.UserCID,
	}
	return models.CreateModelInstance(&generalAccount)
}

func Login(l LoginRequest) (string, error) {
	user, err := models.GetUserByUsername(l.Username)
	if err != nil {
		return "", err
	}
	if user.Password != customHash(l.Password) {
		return "", errors.New("wrong password, try again")
	}

	token, err := createToken(
		RequestUser{
			UserID:   uint(user.ID),
			Username: user.Username,
			UserCID:  user.UserCID,
		},
	)
	if err != nil {
		return "", err
	}
	return token, nil
}
