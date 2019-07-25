package app

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
)


type App struct {
}

func (app *App) Init() error {

	if err := initDataPath(); err != nil {
		return err
	}

	cleanFolder(GetLogPath())
	return nil
}

func removeFile(file string)  {
	err := os.Remove(file)
	if err != nil {
		fmt.Println("file remove Error!")
		fmt.Printf("%s", err)
	} else {
		fmt.Print("file remove OK!")
	}
}

func CleanData()  {
	cleanFolder(GetDataPath())
}

func cleanFolder(folder string){
	files, _ := ioutil.ReadDir(folder)
	for _, file := range files {
		if file.IsDir() {
			continue
		} else {
			fmt.Println("remove file :",file.Name())
			removeFile(folder+file.Name())
		}
	}
}
func initDataPath() (err error) {
	if home, err := Home(); err != nil {
		return fmt.Errorf("Current operating system is not supported，err=[%v] ", err)
	} else {
		switch runtime.GOOS {
		case "darwin":
			app_home_path = home + "/Library/pullup"
			app_keystore_path = app_home_path + "/keystore"
			app_data_path = app_home_path + "/data/"
			app_log_path = app_home_path + "/log/"
			app_config_path = app_home_path + "/config"
			osType = "mac"
			break
		case "windows":
			app_home_path = home + `\AppData\Roaming\pullup`
			app_keystore_path = app_home_path + "\\keystore"
			app_data_path = app_home_path + "\\data\\"
			app_log_path = app_home_path + "\\log\\"
			app_config_path = app_home_path + "\\config"
			osType = "win"
			break
		case "linux":
			app_home_path = home + "/.config/pullup"
			app_keystore_path = app_home_path + "/keystore"
			app_data_path = app_home_path + "/data/"
			app_log_path = app_home_path + "/log/"
			app_config_path = app_home_path + "/config"
			osType = "linux"
			break
		}
	}
	if app_home_path == "" {
		return fmt.Errorf("Current operating system is not supported ")
	}
	subdirectory := []string{app_keystore_path, app_data_path, app_log_path, app_config_path}

	if _, err := os.Stat(app_home_path); os.IsNotExist(err) {
		if err = os.MkdirAll(app_home_path, os.ModePerm); err != nil {
			return fmt.Errorf("Application folder initialization failed，err=[%v] ", err)
		}
		for _, dir := range subdirectory {
			if err = os.MkdirAll(dir, os.ModePerm); err != nil {
				return fmt.Errorf("Application folder initialization failed，err=[%v] ", err)
			}
		}
	} else {
		for _, dir := range subdirectory {
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				if err = os.MkdirAll(dir, os.ModePerm); err != nil {
					return fmt.Errorf("Application folder initialization failed，err=[%v] ", err)
				}
			}
		}
	}
	return nil
}

func GetPath(folder string) string {
	return app_home_path + folder
}

func GetLogPath() string {
	return app_log_path
}

func GetKeystorePath() string {
	return app_keystore_path
}

func GetDataPath() string {
	return app_data_path
}

func GetConfigPath() string {
	return app_config_path
}

func (app *App) GetHomePath() string {
	return app_home_path
}

// Home returns the home directory for the executing user.
//
// This uses an OS-specific method for discovering the home directory.
// An error is returned if a home directory cannot be detected.
func Home() (string, error) {
	user, err := user.Current()
	if nil == err {
		return user.HomeDir, nil
	}
	// cross compile support
	if "windows" == runtime.GOOS {
		return homeWindows()
	}

	// Unix-like system, so just assume Unix
	return homeUnix()
}

func homeUnix() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	// If that fails, try the shell
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}

func homeWindows() (string, error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}

	return home, nil
}
