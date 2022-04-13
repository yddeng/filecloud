package service

import "time"

type authHandler struct {
}

var (
	accessToken       string
	accessTokenExpire time.Time
)

func (*authHandler) login(wait *WaitConn, req struct {
	Username string `json:"username"`
	Password string `json:"password"`
}) {
	logger.Infof("%s %v", wait.GetRoute(), req)
	defer func() { wait.Done() }()

	if req.Username != config.Username || req.Password != config.Password {
		wait.SetResult("用户或密码错误", nil)
		return
	}

	now := time.Now()
	if accessToken == "" || now.After(accessTokenExpire) {
		accessToken = GenToken(20)
		accessTokenExpire = now.Add(time.Hour * 8)
	}

	wait.SetResult("", struct {
		Token string `json:"token"`
	}{Token: accessToken})
}
