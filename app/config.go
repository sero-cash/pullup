package app

//config
var (
	version   = "v0.1.5"
	cleanData = false

	maxUint64  = ^uint64(0)
	fetchCount = uint64(50000)

	// for Asia
	defaultRpcHost = "http://39.98.253.20:8546"
	defaultWebHost = "http://129.211.98.114:3006/web/v0_1_5/"

	//defaultRpcHost = "http://52.199.145.159:8545"
	//defaultWebHost = "http://pullup.sero.cash/v0_1_5/"

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
