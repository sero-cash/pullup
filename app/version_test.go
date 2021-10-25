package app

import (
	"testing"
)

func TestDeCompress(t *testing.T) {

	err :=DeCompressByPath("/Users/tom/Downloads/docs.zip","/Users/tom/Downloads/docs/")
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("DeCompress success")
}