package test

import (
	"meiwobuxing"
	"testing"
)

func TestCreateJT(t *testing.T) {
	u := meiwobuxing.User{Username: "hello", Password: "123333"}
	token, err := u.CreateToken()
	if err != nil {
		println(err.Error())
		return
	}
	println(token)

	println(meiwobuxing.CheckToken(token))
}

func TestTJwt(t *testing.T) {
	u := meiwobuxing.User{Username: "caicai", Password: "159463"}

	login, err := u.Login()
	if err != nil {
		println(err.Error())
		return
	}
	println(login)

	println(meiwobuxing.CheckToken(login))

}