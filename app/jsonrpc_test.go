package app

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestDoRequest(t *testing.T) {
	out, err := DoRequest("http://localhost:2345/web/dapp.json")
	if err != nil {
		return
	}
	dapp := Dapp{}
	err = json.Unmarshal(out, &dapp)
	if err != nil {
		fmt.Println("dapp.json格式不对或者文件不存在. ")
		return
	}
	if dapp.URL != "" && dapp.Author != "" && dapp.Desc !="" && dapp.Img != "" && dapp.Title != "" {
		fmt.Println("Add App to database. ")
	}else {
		fmt.Println("dapp.json格式不正确或者为空 ")
	}
	fmt.Println(string(out))
}
