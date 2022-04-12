package service

import (
	"fmt"
	"github.com/yddeng/dnet/dhttp"
	"testing"
)

func TestMove(t *testing.T) {
	elem := map[string]interface{}{
		"source": []string{"cloud/222", "cloud/test.txt"},
		"target": "cloud/111",
		"move":   true,
	}

	req, _ := dhttp.PostJson("http://127.0.0.1:9987/file/mvcp", elem)
	ret, err := req.ToString()
	fmt.Println(ret, err)
}

func TestRename(t *testing.T) {
	elem := map[string]interface{}{
		"path":    "cloud/111",
		"oldName": "test2.txt",
		"newName": "test.txt",
	}

	req, _ := dhttp.PostJson("http://127.0.0.1:9987/file/rename", elem)
	ret, err := req.ToString()
	fmt.Println(ret, err)
}

func TestList(t *testing.T) {
	elem := map[string]interface{}{
		"path": "cloud",
	}

	req, _ := dhttp.PostJson("http://127.0.0.1:9987/file/list", elem)
	ret, err := req.ToString()
	fmt.Println(ret, err)
}
