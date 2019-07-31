package app

//config
var (
	version   = "v0.1.4"
	cleanData = true

	maxUint64  = ^uint64(0)
	fetchCount = uint64(50000)

	defaultRpcHost = "http://39.98.253.20:8546"
	//defaultWebHost = "http://129.211.98.114:3006"
	defaultWebHost = "http://127.0.0.1:2345"
	osType         = ""
	rpcHost        = ""
	webHost        = ""

	app_home_path     string
	app_keystore_path string
	app_log_path      string
	app_config_path   string
	app_data_path     string
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
