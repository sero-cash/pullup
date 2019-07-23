package main

import (
	"encoding/json"
	"flag"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/sero-cash/go-czero-import/cpt"
	"github.com/sero-cash/go-sero/pullup/app"
	"github.com/sero-cash/go-sero/pullup/common/logex"
	"github.com/sero-cash/go-sero/pullup/common/transport"
	"github.com/sero-cash/go-sero/pullup/lorca"
	"net/http"
	"runtime"
	"time"
)

func main() {

	cpt.ZeroInit_OnlyInOuts()

	rpcHostCustomer := flag.String("rpcHost","","--rpcHost set rpc host")
	webHostCustomer := flag.String("webHost","","--webHost set web host")
	flag.Parse()

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

	// init sero light client
	app.NewSeroLight()

	//banding http handle
	privateAccountApi := app.NewServiceAPI()

	privateAccountApi.InitHost(*rpcHostCustomer,*webHostCustomer)

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

	http.HandleFunc("/web/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	http.HandleFunc("/rpc", func(w http.ResponseWriter, r *http.Request) {
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

	// init ui
	args := []string{}
	if runtime.GOOS == "linux" {
		args = append(args, "--class=Lorca")
	}
	ui, err := lorca.New("", "", 1400, 768, args...)
	if err != nil {
		logex.Fatal(err)
	}
	defer ui.Close()

	go func() {
		time.Sleep(time.Second*1)
		// A simple way to know when UI is ready (uses body.onload event in JS)
		if err = ui.Bind("start", func() {
			logex.Info("UI is ready")
		}); err != nil {
			logex.Fatal(err)
		}
		if err = ui.Load(app.GetWebHost()+"/web/"); err != nil {
			logex.Fatal(err)
		}
	}()

	err = http.ListenAndServe(":2345", nil)
	if err != nil {
		logex.Fatal(err)
	}
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}

type rpcParams struct {
	Method string `json:"method"`
	Params interface{} `json:"params"`
}