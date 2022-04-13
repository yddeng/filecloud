package service

import (
	"fmt"
	"github.com/yddeng/dnet/dhttp"
	"testing"
)

var (
	address = "127.0.0.1:9987"
)

func authLogin(t *testing.T, Username, Password string) string {
	req, _ := dhttp.PostJson(fmt.Sprintf("http://%s/auth/login", address), struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{Username: Username, Password: Password})

	ret, err := req.ToString()
	if err != nil {
		t.Fatal(err)
	}
	return ret
}

func TestAuthLogin(t *testing.T) {
	ret := authLogin(t, "yddeng", "123456")
	fmt.Println(ret)
}
