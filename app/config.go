package app

import (
	"strings"
	"time"
)

//config
var (
	version   = "v0.2.2"
	cleanData = false

	maxUint64  = ^uint64(0)
	fetchCount = uint64(500000)

	remoteConfig = ""

	versionUrl = ""

	osType  = ""
	rpcHost = ""
	webHost = ""

	app_home_path     string
	app_keystore_path string
	app_log_path      string
	app_config_path   string
	app_data_path     string
	app_cache_path    string

	IsDev = false
	CmdPath = ""

	//default setting
	useZNum = uint64(1958696)

	RemoteVersion = TVersion{}
)

func SetVersionUrl(url string)  {
	versionUrl = url
}

func GetVersionUrl() string  {
	return versionUrl
}

func SetRemoteConfig(config string)  {
	remoteConfig = config
}

func GetRemoteConfig() string {
	return remoteConfig
}
func GetVersion() string {
	return version
}

func setRpcHost(s string) {
	rpcHost = s
}
func GetRpcHost() string {
	return rpcHost
}

func setWebHost(s string) {
	webHost = s
}
func GetWebHost() string {
	return localDocs
}
func GetOsType() string {
	return osType
}

type Node struct {
	Id      string `json:"id"`
	Network string `json:"network"`
	Name    string `json:"name"`
	Rpc     string `json:"rpc"`
	Web     string `json:"web"`
}

type RpcConfig struct {
	Default Node   `json:"default"`
	Host    []Node `json:"host"`
}

func IsZH() bool {
	location := time.Now().String()
	return strings.Index(location, "+0800") > -1
}