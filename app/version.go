package app

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/sero-cash/go-sero/pullup/common/logex"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"time"
)

var (
	localDocs = "http://127.0.0.1:3646/docs"
)

func CheckVersion() {
	for {
		localVersion := TVersion{}
		err := httpGet(localDocs+"/version.json", &localVersion)
		if err != nil {
			logex.Error(err.Error())
			time.Sleep(60 * time.Second)
			continue
		}
		remoteVersion := TVersion{}
		err =httpGet(GetVersionUrl(), &remoteVersion)
		if err != nil{
			logex.Error(err.Error())
			time.Sleep(60 * time.Second)
			continue
		}

		//store remoteVersion mem
		if remoteVersion.Version.App != ""{
			RemoteVersion = remoteVersion
		}
		
		if localVersion.Version.Docs != remoteVersion.Version.Docs {
			logex.Infof("Download file=[docs.zip] begin , version = [%s]", remoteVersion.Version.Docs)
			downDocsUrl := remoteVersion.Version.DocsUrl["en"]
			if IsZH() {
				downDocsUrl= remoteVersion.Version.DocsUrl["zh"]
			}
			targetZip := CmdPath + "/docs.zip"
			webPath := CmdPath+"/docs/"
			switch runtime.GOOS {
			case "windows":
				targetZip = CmdPath + "\\docs.zip"
				webPath = CmdPath + "\\docs\\"
				break
			}

			fmt.Println(downDocsUrl)
			res, err := http.Get(downDocsUrl)
			if err != nil {
				logex.Error(err.Error())
				time.Sleep(60 * time.Second)
				continue
			}
			f, err := os.Create(targetZip)
			if err != nil {
				logex.Error(err.Error())
				time.Sleep(60 * time.Second)
				continue
			}
			io.Copy(f, res.Body)
			logex.Infof("Download file=[docs.zip] success , version = [%s]", remoteVersion.Version.Docs)
			logex.Info("DeCompress file=[docs.zip] begin")
			if err := DeCompressByPath(targetZip, webPath); err != nil {
				logex.Errorf("DeCompress file=[docs.zip], error=[%s]", err.Error())
				time.Sleep(60 * time.Second)
				continue
			}
			logex.Info("DeCompress file=[docs.zip] success")
		}else{
			logex.Infof("No new docs update local version=[%s], waiting... 60s for check update.",localVersion.Version.Docs)
		}
		time.Sleep(60 * time.Second)
	}
}

type TVersion struct {
	Version     VersionN `json:"version"`
	Description map[string][]string `json:"description"`
}

type VersionN struct {
	Docs    string `json:"docs"`
	DocsUrl map[string]string `json:"docsUrl"`
	App     string `json:"app"`
	AppUrl  map[string]string `json:"appUrl"`
}

func httpGet(url string, rest interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		logex.Errorf("Get url=[%s] error=[%s] ", url, err.Error())
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logex.Errorf("ReadAll url=[%s] error=[%s] ", url, err.Error())
		return err
	}
	return json.Unmarshal(body, rest)
}

/**
@tarFile：压缩文件路径
@dest：解压文件夹
*/
func DeCompressByPath(tarFile, dest string) error {
	srcFile, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	return DeCompress(srcFile, dest)
}

/**
@zipFile：压缩文件
@dest：解压之后文件保存路径
*/
func DeCompress(srcFile *os.File, dest string) error {
	zipFile, err := zip.OpenReader(srcFile.Name())
	if err != nil {
		logs.Error("Unzip File Error：", err.Error())
		return err
	}
	defer zipFile.Close()
	for _, innerFile := range zipFile.File {
		info := innerFile.FileInfo()
		if info.IsDir() {
			err = os.MkdirAll(dest+innerFile.Name, os.ModePerm)
			if err != nil {
				logs.Error("Unzip File Error : " + err.Error())
				return err
			}
			continue
		}
		srcFile, err := innerFile.Open()
		if err != nil {
			logs.Error("Unzip File Error : " + err.Error())
			continue
		}
		defer srcFile.Close()
		newFile, err := os.Create(dest + innerFile.Name)
		if err != nil {
			logs.Error("Unzip File Error : " + err.Error())
			continue
		}
		io.Copy(newFile, srcFile)
		newFile.Close()
	}
	return nil
}
