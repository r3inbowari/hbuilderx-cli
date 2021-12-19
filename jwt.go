package meiwobuxing

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	. "github.com/r3inbowari/zlog"
	"github.com/sirupsen/logrus"
	"time"
)

const salt = "&rbc2gx4-+="

// CreateToken create token
func (u *User) CreateToken() (string, error) {
	claims := &jwt.StandardClaims{
		NotBefore: time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(GetConfig(false).JwtTimeout)).Unix(),
		Issuer:    "caicai",
		Id:        u.Username,
	}
	Log.WithFields(logrus.Fields{"uid": u.Username}).Info("[JWT] create token")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(GetConfig(false).JwtSecret))
	if err != nil {
		return "", err
	}
	return ss, nil
}

func (u *User) Login() (string, error) {

	if GetMD5WithSalt(fmt.Sprintf("%s%s", u.Username, u.Password), salt) == GetConfig(false).JwtMD5 {
		return u.CreateToken()
	}
	return "", errors.New("401")
}

// CheckToken check token
func CheckToken(token string) error {
	_, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return []byte(GetConfig(false).JwtSecret), nil
	})
	return err
}

func CalcJwtMD5(username, password string) string {
	return GetMD5WithSalt(fmt.Sprintf("%s%s", username, password), salt)
}
