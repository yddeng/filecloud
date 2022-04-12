package service

import (
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/yddeng/dnet/dhttp"
	"testing"
)

func TestShare(t *testing.T) {
	elem := map[string]interface{}{
		"path":     "cloud",
		"filename": []string{"111", "babel.config.js"},
		"deadline": 0,
	}

	req, _ := dhttp.PostJson("http://127.0.0.1:9987/shared/create", elem)
	ret, err := req.ToString()
	fmt.Println(ret, err)

	elem = map[string]interface{}{
		"path":  "cloud",
		"token": gjson.Get(ret, "data.token").String(),
	}

	req, _ = dhttp.PostJson(gjson.Get(ret, "data.route").String(), elem)
	ret, err = req.ToString()
	fmt.Println(ret, err)
}
