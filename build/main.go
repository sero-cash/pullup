package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var dir string
func main() {

	var err error
	dir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("dir：",dir)

	//dir = "/Users/huangw/Documents/codes/go/src/github.com/sero-cash/go-sero/light-wallet/build"

	command := ""
	switch runtime.GOOS {
	case "darwin":
		command = dir+"/bin/light-wallet-darwin-10.6-amd64"
		break
	case "windows":
		command = dir+"\\bin\\light-wallet-windows-4.0-amd64.exe"
		break
	case "linux":
		command = dir+"/bin/light-wallet-linux-4.0-amd64"
		break
	}
	if command != "" {
		execCommand(command, nil)
	}else{
		log.Fatal("This OS is not supported,exit!")
	}
}

func ChangeYourCmdEnvironment(cmd *exec.Cmd) error {
	env := os.Environ()
	cmdEnv := []string{}

	for _, e := range env {
		i := strings.Index(e, "=")
		if i > 0 && (e[:i] == "ENV_NAME") {
		} else {
			var tempstring = e
			if e[:i] =="Path"{
				tempstring=tempstring + ";"+ dir+"\\czero\\lib\\"
			}
			cmdEnv = append(cmdEnv, tempstring)
		}
	}
	switch runtime.GOOS {
	case "darwin":
		fmt.Println("darwin:","DYLD_LIBRARY_PATH="+dir+"/czero/lib/")
		cmdEnv = append(cmdEnv, "DYLD_LIBRARY_PATH="+dir+"/czero/lib/")
		break
	case "windows":
		fmt.Println("windows:","Path="+dir+"\\czero\\lib\\")
		//cmdEnv = append(cmdEnv, "Path="+dir+"\\czero\\lib\\")
		break
	case "linux":
		cmdEnv = append(cmdEnv, "LD_LIBRARY_PATH="+dir+"/czero/lib/")
		break
	}
	cmd.Env = cmdEnv
	return nil
}

func execCommand(commandName string, params []string) bool {
	cmd := exec.Command(commandName, params...)
	ChangeYourCmdEnvironment(cmd)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return false
	}
	cmd.Start()
	reader := bufio.NewReader(stdout)
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			fmt.Errorf("app exit,err:%s ", err2.Error())
			break
		}
		fmt.Println(line)
	}
	cmd.Wait()
	return true
}
