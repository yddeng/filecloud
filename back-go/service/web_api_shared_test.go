package service

import (
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/yddeng/dnet/dhttp"
	"strings"
	"testing"
)

func TestShareInfo(t *testing.T) {
	ret := authLogin(t, "yddeng", "123456")
	fmt.Println(ret, gjson.Get(ret, "data.token").String())

	token := gjson.Get(ret, "data.token").String()

	elem := map[string]interface{}{
		"path":     "cloud",
		"filename": []string{"111", "babel.config.js"},
		"deadline": 7,
	}

	req, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/shared/create", address), "POST")
	req.SetHeader("Access-Token", token)
	req.WriteJSON(elem)
	ret, err := req.ToString()
	fmt.Println(ret, err)

	elem = map[string]interface{}{
		"key":         strings.TrimPrefix(gjson.Get(ret, "data.route").String(), "http://127.0.0.1:9987/shared/s/"),
		"sharedToken": gjson.Get(ret, "data.sharedToken").String(),
	}

	req, _ = dhttp.NewRequest(fmt.Sprintf("http://%s/shared/info", address), "POST")
	//req.SetHeader("Access-Token", token)
	req.WriteJSON(elem)
	ret, err = req.ToString()
	fmt.Println(ret, err)
}

func TestShareList(t *testing.T) {
	// 链接：http://127.0.0.1:9987/shared/s/eHr6Qp2ji9mqKmoN  提取码：F9bS
	elem := map[string]interface{}{
		"key":         "6JvwhI0a6zddfZqW",
		"path":        "cloud",
		"sharedToken": "PlA7",
	}

	req, _ := dhttp.PostJson(fmt.Sprintf("http://%s/shared/list", address), elem)
	ret, err := req.ToString()
	fmt.Println(ret, err)
}
