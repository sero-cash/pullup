package app

//config
var (
	version   = "v0.1.6"
	cleanData = true

	maxUint64  = ^uint64(0)
	fetchCount = uint64(50000)

	// for Asia
	remoteRpcHost  = "https://sero-media-1256272584.cos.ap-shanghai.myqcloud.com/pullup/config/v0.1.6/node.json"

	//defaultRpcHost = "http://52.199.145.159:8545"
	//defaultWebHost = "http://pullup.sero.cash/v0_1_6/"
	//remoteRpcHost  = "https://sero-media-1256272584.cos.ap-shanghai.myqcloud.com/pullup/config/v0.1.6/node.json"

	osType  = ""
	rpcHost = ""
	webHost = ""

	app_home_path     string
	app_keystore_path string
	app_log_path      string
	app_config_path   string
	app_data_path     string

	IsDev = false
)

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
