package app

//config
var (
	version   = "v0.1.8"
	cleanData = true

	maxUint64  = ^uint64(0)
	fetchCount = uint64(50000)

	remoteConfig = ""

	//remoteRpcHost  = "https://sero-media.s3-ap-southeast-1.amazonaws.com/clients/node-global.json"
	//defaultRpcHost = "http://52.199.145.159:8545"
	//defaultWebHost = "http://pullup.sero.cash/v0_1_6/"

	osType  = ""
	rpcHost = ""
	webHost = ""

	app_home_path     string
	app_keystore_path string
	app_log_path      string
	app_config_path   string
	app_data_path     string

	IsDev = false
	CmdPath = ""

	//default setting
	useZNum = uint64(100)
)

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
	return webHost
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
