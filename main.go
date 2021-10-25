package main

import (
	"encoding/json"
	"flag"
	"fmt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/sero-cash/go-czero-import/superzk"
	"github.com/sero-cash/go-sero/pullup/app"
	"github.com/sero-cash/go-sero/pullup/common/logex"
	"github.com/sero-cash/go-sero/pullup/common/transport"
	"github.com/sero-cash/go-sero/pullup/lorca"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	h bool
	g string
	w string
	d bool
	c string
	l string
)

var (
	crossOrigin = "http://127.0.0.1:3646"
)

func init() {
	flag.BoolVar(&h, "h", false, "this help")
	flag.StringVar(&g, "g", "", "set `gero rpc host`")
	flag.StringVar(&w, "w", "", "set `web url`")
	flag.BoolVar(&d, "d", false, "set `dev env`")
	flag.StringVar(&c, "c", "", "set `component path`")
	flag.StringVar(&l, "l", "zh", "set `location`")
	flag.Usage = usage
}
func usage() {
	fmt.Fprintf(os.Stderr, `pullup
Usage: pullup [-g gero] [-w web]   [-d dev] [-c component] [-l local]
Options:
`)
	flag.PrintDefaults()
}

func main() {

	flag.Parse()
	if h {
		usage()

		return
	}
	log.Println("isDev",d)
	//set start flag
	app.IsDev = d
	app.CmdPath = c

	// Setting global env
	lightApp := app.App{}
	if err := lightApp.Init(); err != nil {
		panic(err)
	}

	// setup log
	log := logex.Log{Name: "pullup", Path: app.GetLogPath()}
	if logFile, err := log.Setup(); err != nil {
		panic(err)
	} else {
		defer logFile.Close()
	}

	if exe := lorca.ChromeExecutable(); exe == "" {
		logex.Info("No Chrome ,go to download ")
		lorca.PromptDownload()
		return
	}

	//init Zero import
	superzk.ZeroInit_OnlyInOuts()
	logex.Info("ZeroInit_OnlyInOuts successful! ")

	if app.IsDev {
		crossOrigin = "*"
	}

	setCmdPath()


	fmt.Println("app.CmdPath:",app.CmdPath)
	go func() {
		h := http.FileServer(http.Dir(app.CmdPath + "/docs"))
		http.Handle("/docs/", http.StripPrefix("/docs/", h)) // 启动静态文件服务
		//Header().Set("Expires", time.Now().Format("MON, 02 Jan 2006 15:04:05 GMT"))
		http.ListenAndServe(":3646", nil)
	}()
	// init sero light client
	app.NewSeroLight()
	logex.Info("NewSeroLight successful! ")

	registerHttpHandler()

	// start up a http server
	ln, err := net.Listen("tcp", "127.0.0.1:2345")
	if err != nil {
		logex.Fatal(err)
	}
	defer ln.Close()
	go func() {
		// Set up your http server here
		logex.Fatal(http.Serve(ln, nil))
	}()

	if !app.IsDev {
		go app.CheckVersion()
	}


	// init ui
	args := []string{"--disable-backgrounding-occluded-windows"}
	if runtime.GOOS == "linux" {
		args = append(args, "--class=Lorca")
	}

	ui, err := lorca.New(app.GetWebHost(), app.GetCachePath(), 1400, 768, args...)
	if err != nil {
		logex.Fatal(err)
	}
	defer ui.Close()
	go func() {
		// A simple way to know when UI is ready (uses body.onload event in JS)
		if err = ui.Bind("start", func() {
			logex.Info("UI is ready")
		}); err != nil {
			logex.Fatal(err)
		}
		if err = ui.Load(app.GetWebHost()); err != nil {
			logex.Fatal(err)
		}

	}()
	<-ui.Done()

}

func setCmdPath() {
	var err error
	app.CmdPath, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logex.Fatal(err.Error())
	}
	logex.Info("Cmd Path >>> ", app.CmdPath)
}



func registerHttpHandler() {
	if app.IsZH() {
		app.SetRemoteConfig("http://127.0.0.1:3646/docs/node-asia.json")
		app.SetVersionUrl("https://pullup.sero.cash/version.json")
	} else {
		app.SetRemoteConfig("http://127.0.0.1:3646/docs/node-global.json")
		app.SetVersionUrl("https://pullup.sero.cash/version.json")
	}

	//banding http handle
	privateAccountApi := app.NewServiceAPI()
	privateAccountApi.InitHost(g)

	createAccountHandler := httptransport.NewServer(
		app.MakeAccountCreateEndpoint(privateAccountApi),
		transport.DecodeRequest,
		transport.EncodeResponse,
	)
	http.Handle("/account/create", accessControl(createAccountHandler))
	//upload keystore
	http.HandleFunc("/account/import/keystore", privateAccountApi.UploadKeystoreHandler())
	importWithMnemonicHandler := httptransport.NewServer(
		app.MakeAccountImportWithMnemonicEndpoint(privateAccountApi),
		transport.DecodeRequest,
		transport.EncodeResponse,
	)
	http.Handle("/account/import/mnemonic", accessControl(importWithMnemonicHandler))
	importWithPrivateHandler := httptransport.NewServer(
		app.MakeAccountImportWithPrivateKeyEndpoint(privateAccountApi),
		transport.DecodeRequest,
		transport.EncodeResponse,
	)
	http.Handle("/account/import/private", accessControl(importWithPrivateHandler))
	exportMnemonicHandler := httptransport.NewServer(
		app.MakeAccountExportMnemonicEndpoint(privateAccountApi),
		transport.DecodeRequest,
		transport.EncodeResponse,
	)
	http.Handle("/account/export/mnemonic", accessControl(exportMnemonicHandler))
	accountListHandler := httptransport.NewServer(
		app.MakeAccountListEndpoint(privateAccountApi),
		transport.DecodeRequest,
		transport.EncodeResponse,
	)
	http.Handle("/account/list", accessControl(accountListHandler))
	accountDetailHandler := httptransport.NewServer(
		app.MakeAccountDetailEndpoint(privateAccountApi),
		transport.DecodeRequest,
		transport.EncodeResponse,
	)
	http.Handle("/account/detail", accessControl(accountDetailHandler))
	accountBalanceHandler := httptransport.NewServer(
		app.MakeAccountBalanceEndpoint(privateAccountApi),
		transport.DecodeRequest,
		transport.EncodeResponse,
	)
	http.Handle("/account/balance", accessControl(accountBalanceHandler))
	txListHandler := httptransport.NewServer(
		app.MakeTxListEndpoint(privateAccountApi),
		transport.DecodeRequest,
		transport.EncodeResponse,
	)
	http.Handle("/tx/list", accessControl(txListHandler))
	txNumHandler := httptransport.NewServer(
		app.MakeTxNumEndpoint(privateAccountApi),
		transport.DecodeRequest,
		transport.EncodeResponse,
	)
	http.Handle("/tx/num", accessControl(txNumHandler))
	txTransferHandler := httptransport.NewServer(
		app.MakeTxSendEndpoint(privateAccountApi),
		transport.DecodeRequest,
		transport.EncodeResponse,
	)
	http.Handle("/tx/transfer", accessControl(txTransferHandler))

	keyPathandler := httptransport.NewServer(
		app.MakeDataPathEndpoint(privateAccountApi),
		transport.DecodeRequest,
		transport.EncodeResponse,
	)
	http.Handle("/path/keystore", accessControl(keyPathandler))

	addOutHandler := httptransport.NewServer(
		app.MakeAddOutEndpoint(privateAccountApi),
		transport.DecodeRequest,
		transport.EncodeResponse,
	)
	http.Handle("/tx/addOut", accessControl(addOutHandler))

	decimalHandler := httptransport.NewServer(
		app.MakeCurrencyDecimalEndpoint(privateAccountApi),
		transport.DecodeRequest,
		transport.EncodeResponse,
	)
	http.Handle("/decimal", accessControl(decimalHandler))
	stakePoolListHandler := httptransport.NewServer(
		app.MakeStakePoolEndpoint(privateAccountApi),
		transport.DecodeRequest,
		transport.EncodeResponse,
	)
	http.Handle("/stake", accessControl(stakePoolListHandler))
	registerPoolhandler := httptransport.NewServer(
		app.MakeRegisterShareEndpoint(privateAccountApi),
		transport.DecodeRequest,
		transport.EncodeResponse,
	)
	http.Handle("/stake/register", accessControl(registerPoolhandler))
	closePoolhandler := httptransport.NewServer(
		app.MakeCloseShareEndpoint(privateAccountApi),
		transport.DecodeRequest,
		transport.EncodeResponse,
	)
	http.Handle("/stake/close", accessControl(closePoolhandler))
	buySharehandler := httptransport.NewServer(
		app.MakeBuyShareEndpoint(privateAccountApi),
		transport.DecodeRequest,
		transport.EncodeResponse,
	)
	http.Handle("/stake/buyShare", accessControl(buySharehandler))
	getSharehandler := httptransport.NewServer(
		app.MakeGetShareEndpoint(privateAccountApi),
		transport.DecodeRequest,
		transport.EncodeResponse,
	)
	http.Handle("/stake/getShare", accessControl(getSharehandler))
	getTxReceipthandler := httptransport.NewServer(
		app.MakeGetTransactionReceiptEndpoint(privateAccountApi),
		transport.DecodeRequest,
		transport.EncodeResponse,
	)
	http.Handle("/tx/getTxReceipt", accessControl(getTxReceipthandler))
	getBlockNumberhandler := httptransport.NewServer(
		app.MakeGetBlockNumberEndpoint(privateAccountApi),
		transport.DecodeRequest,
		transport.EncodeResponse,
	)
	http.Handle("/sero/getBlockNumber", accessControl(getBlockNumberhandler))
	getMySharesHandler := httptransport.NewServer(
		app.MakeGetMySharesEndpoint(privateAccountApi),
		transport.DecodeRequest,
		transport.EncodeResponse,
	)
	http.Handle("/share/my", accessControl(getMySharesHandler))
	changeNetworkHandler := httptransport.NewServer(
		app.MakeChangeNetworkEndpoint(privateAccountApi),
		transport.DecodeRequest,
		transport.EncodeResponse,
	)
	http.Handle("/network/change", accessControl(changeNetworkHandler))
	openFileHandler := httptransport.NewServer(
		app.MakeOpenFileEndpoint(privateAccountApi),
		transport.DecodeRequest,
		transport.EncodeResponse,
	)
	http.Handle("/file/open", accessControl(openFileHandler))
	dappHandler := httptransport.NewServer(
		app.MakeSetDappsEndpoint(privateAccountApi),
		transport.DecodeRequest,
		transport.EncodeResponse,
	)
	http.Handle("/dapp/set", accessControl(dappHandler))
	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
		if r.Method == "OPTIONS" {
			return
		}
		json.NewEncoder(w).Encode(app.GetVersion())
		return
	})

	http.HandleFunc("/remoteVersion", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
		if r.Method == "OPTIONS" {
			return
		}
		json.NewEncoder(w).Encode(app.RemoteVersion)
		return
	})

	http.HandleFunc("/rpc", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
		if r.Method == "OPTIONS" {
			return
		}
		rpcParams := rpcParams{}
		if err := json.NewDecoder(r.Body).Decode(&rpcParams); err != nil {
			json.NewEncoder(w).Encode(err.Error())
			return
		}
		sync := app.Sync{RpcHost: app.GetRpcHost(), Method: rpcParams.Method, Params: rpcParams.Params}
		jsonResp, err := sync.Do()
		if err != nil {
			json.NewEncoder(w).Encode(err.Error())
			return
		}
		json.NewEncoder(w).Encode(jsonResp)
		return
	})
	//pre collect data
	http.HandleFunc("/pullup_rpc", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
		if r.Method == "OPTIONS" {
			return
		}
		req := app.JSONRpcReq{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			json.NewEncoder(w).Encode(err.Error())
			return
		}
		resp := app.HandlePullupRpc(req, privateAccountApi)
		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			logex.Errorf("HandlePullupRpc resp json err: ", err)
		}
		return
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
		if r.Method == "OPTIONS" {
			return
		}
		req := app.JSONRpcReq{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			json.NewEncoder(w).Encode(err.Error())
			return
		}
		if strings.Split(string(req.Method), "_")[0] == "pullup" {
			resp := app.HandlePullupRpc(req, privateAccountApi)
			err := json.NewEncoder(w).Encode(resp)
			if err != nil {
				logex.Errorf("HandlePullupRpc resp json err: ", err)
			}
		} else {
			sync := app.Sync{RpcHost: app.GetRpcHost(), Method: string(req.Method), Params: req.Params}
			jsonResp, err := sync.Do()
			if err != nil {
				json.NewEncoder(w).Encode(err.Error())
				return
			}
			json.NewEncoder(w).Encode(jsonResp)
		}
		return
	})
	logex.Info("http handler loaded successful.")
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", crossOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}

type rpcParams struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}
