package service

import (
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/yddeng/dnet/dhttp"
	"strings"
	"testing"
)

func TestShare(t *testing.T) {
	elem := map[string]interface{}{
		"path":     "cloud/111",
		"filename": []string{"222", "babel.config.js"},
		"deadline": 0,
	}

	req, _ := dhttp.PostJson("http://127.0.0.1:9987/shared/create", elem)
	ret, err := req.ToString()
	fmt.Println(ret, err)

	elem = map[string]interface{}{
		"key":         strings.TrimPrefix(gjson.Get(ret, "data.route").String(), "http://127.0.0.1:9987/shared/s/"),
		"path":        "cloud/111",
		"sharedToken": gjson.Get(ret, "data.sharedToken").String(),
	}

	req, _ = dhttp.PostJson("http://127.0.0.1:9987/shared/list", elem)
	ret, err = req.ToString()
	fmt.Println(ret, err)
}

func TestShareList(t *testing.T) {
	// 链接：http://127.0.0.1:9987/shared/s/eHr6Qp2ji9mqKmoN  提取码：F9bS
	elem := map[string]interface{}{
		"key":         "eHr6Qp2ji9mqKmoN",
		"path":        "cloud",
		"sharedToken": "F9bS",
	}

	req, _ := dhttp.PostJson("http://127.0.0.1:9987/shared/list", elem)
	ret, err := req.ToString()
	fmt.Println(ret, err)
}
