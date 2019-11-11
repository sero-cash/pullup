package app

import (
	"bufio"
	"fmt"
	"github.com/sero-cash/go-sero/pullup/common/logex"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func enterPassword(title string) (password string) {
	password, err := executeWebview("password", "-t", title, "-b", "Confirm")
	if err != nil {
		logex.Error("executeWebview err: ", err.Error())
		return ""
	}
	return password
}

func executeWebview(method string, args ...string) (string, error) {
	var err error
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if CmdPath != "" {
		dir = CmdPath
	}
	if err != nil {
		logex.Error("executeWebview filepath err: ", err.Error())
		return "", err
	}
	command := ""
	switch runtime.GOOS {
	case "darwin":
		if method == "password" {
			command = dir + "/password-darwin-10.6-amd64"
		} else if method == "msgbox" {
			command = dir + "/msgbox-darwin-10.6-amd64"
		}
		break
	case "windows":
		if method == "password" {
			command = dir + `\password-windows-4.0-amd64.exe`
		} else if method == "msgbox" {
			command = dir + `\msgbox-windows-4.0-amd64.exe`
		}
		break
	case "linux":
		command = dir + "/pullup-linux-4.0-amd64"
		break
	}
	if command != "" {
		cmd := exec.Command(command, args...)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Println("exec popu window err: ", err)
			return "", err
		}
		cmd.Start()
		reader := bufio.NewReader(stdout)
		for {
			line, err2 := reader.ReadString('\n')
			if err2 != nil || io.EOF == err2 {
				fmt.Errorf("app exit,err:%s ", err2.Error())
				break
			}
			return strings.TrimSpace(line), nil
		}
		cmd.Wait()
		//
		//out, err := cmd.Output()
		//if err != nil {
		//	logex.Error("executeWebview Command err: ", err.Error())
		//	return "", err
		//}
		return "", fmt.Errorf("Please enter your account password ")
	} else {
		logex.Error("This OS is not supported,exit!")
		return "", err
	}
}
