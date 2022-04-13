package service

import (
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/yddeng/dnet/dhttp"
	"testing"
)

func TestMove(t *testing.T) {
	ret := authLogin(t, "yddeng", "123456")
	t.Log(ret, gjson.Get(ret, "data.token").String())

	token := gjson.Get(ret, "data.token").String()

	elem := map[string]interface{}{
		"source": []string{"cloud/222", "cloud/test.txt"},
		"target": "cloud/111",
		"move":   true,
	}

	req, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/file/mvcp", address), "POST")
	req.SetHeader("Access-Token", token)
	req.WriteJSON(elem)

	ret, err := req.ToString()
	fmt.Println(ret, err)
}

func TestRename(t *testing.T) {
	ret := authLogin(t, "yddeng", "123456")
	t.Log(ret, gjson.Get(ret, "data.token").String())

	token := gjson.Get(ret, "data.token").String()

	elem := map[string]interface{}{
		"path":    "cloud/111",
		"oldName": "test2.txt",
		"newName": "test.txt",
	}

	req, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/file/rename", address), "POST")
	req.SetHeader("Access-Token", token)
	req.WriteJSON(elem)

	ret, err := req.ToString()
	fmt.Println(ret, err)
}

func TestList(t *testing.T) {
	ret := authLogin(t, "yddeng", "123456")
	t.Log(ret, gjson.Get(ret, "data.token").String())

	token := gjson.Get(ret, "data.token").String()

	elem := map[string]interface{}{
		"path": "cloud",
	}

	req, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/file/list", address), "POST")
	req.SetHeader("Access-Token", token)
	req.WriteJSON(elem)
	ret, err := req.ToString()
	fmt.Println(ret, err)
}
