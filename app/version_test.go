package app

import (
	"testing"
)

func TestDeCompress(t *testing.T) {

	err :=DeCompressByPath("/Users/huangw/Downloads/docs.zip","/Users/huangw/Downloads/docs/")
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("DeCompress success")
}