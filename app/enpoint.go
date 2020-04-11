package app

import (
	"context"
	"encoding/json"
	"github.com/sero-cash/go-sero/common/hexutil"
	"os/exec"
	"strconv"
	"sync"

	"github.com/go-kit/kit/endpoint"
	"github.com/sero-cash/go-sero/pullup/common/errorcode"
	"github.com/sero-cash/go-sero/pullup/common/logex"
	"github.com/sero-cash/go-sero/pullup/common/transport"
	"github.com/sero-cash/go-sero/pullup/common/utils"
	"github.com/sero-cash/go-sero/pullup/common/validator"
)

var wg sync.WaitGroup

type AccountCreateReq struct {
	Passphrase string `json:"passphrase"`
}

func MakeAccountCreateEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transport.Request)
		response := transport.Response{}
		response.SetBaseResponseSuccess()

		if ok, err := validator.ValidateBaseRequestParam(req.Base); !ok {
			response.SetBaseResponse(errorcode.InvalidBaseParameters, err.Error())
			return response, nil
		}

		accountCreateReq := AccountCreateReq{}
		utils.Convert(req.Biz, &accountCreateReq)

		if accountCreateReq.Passphrase == "" {
			response.SetBaseResponse(errorcode.FAIL_CODE, "passphrase is nil")
			return response, nil
		}

		//get NODE block  number

		_, blockNumber := getRemoteBlockNumber()

		resp, err := service.NewAccountWithMnemonic(accountCreateReq.Passphrase, uint64(blockNumber))
		if err != nil {
			response.SetBaseResponse(errorcode.FAIL_CODE, err.Error())
		} else {
			response.SetBizResponse(resp)
		}
		return response, nil
	}
}

func getRemoteBlockNumber() (error, hexutil.Uint64) {
	sync := Sync{RpcHost: GetRpcHost(), Method: "sero_blockNumber", Params: []interface{}{}}
	jsonResp, err := sync.Do()
	var blockNumber hexutil.Uint64
	json.Unmarshal(*jsonResp.Result, &blockNumber)
	return err, blockNumber
}

type accountImportWithMnemonicReq struct {
	Mnemonic   string `json:"mnemonic"`
	Passphrase string `json:"passphrase"`
}

func MakeAccountImportWithMnemonicEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transport.Request)
		response := transport.Response{}
		response.SetBaseResponseSuccess()

		if ok, err := validator.ValidateBaseRequestParam(req.Base); !ok {
			response.SetBaseResponse(errorcode.InvalidBaseParameters, err.Error())
			return response, nil
		}
		aimq := accountImportWithMnemonicReq{}
		utils.Convert(req.Biz, &aimq)

		resp, err := service.ImportAccountFromMnemonic(aimq.Mnemonic, aimq.Passphrase)
		if err != nil {
			response.SetBaseResponse(errorcode.FAIL_CODE, err.Error())
		} else {
			response.SetBizResponse(resp)
		}

		return response, nil
	}
}

type accountImportWithPrivateKeyReq struct {
	PrivateKey string `json:"private_key"`
	Passphrase string `json:"passphrase"`
}

func MakeAccountImportWithPrivateKeyEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transport.Request)
		response := transport.Response{}
		response.SetBaseResponseSuccess()

		if ok, err := validator.ValidateBaseRequestParam(req.Base); !ok {
			response.SetBaseResponse(errorcode.InvalidBaseParameters, err.Error())
			return response, nil
		}

		aipq := accountImportWithPrivateKeyReq{}
		utils.Convert(req.Biz, &aipq)

		resp, err := service.ImportAccountFromRawKey(aipq.PrivateKey, aipq.Passphrase, 0, 2)
		if err != nil {
			response.SetBaseResponse(errorcode.FAIL_CODE, err.Error())
		} else {
			response.SetBizResponse(resp)
		}

		return response, nil
	}
}

type accountExportMnemonic struct {
	Passphrase string `json:"passphrase"`
	Address    string `json:"address"`
}

func MakeAccountExportMnemonicEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transport.Request)
		response := transport.Response{}
		response.SetBaseResponseSuccess()

		if ok, err := validator.ValidateBaseRequestParam(req.Base); !ok {
			response.SetBaseResponse(errorcode.InvalidBaseParameters, err.Error())
			return response, nil
		}

		aem := accountExportMnemonic{}
		utils.Convert(req.Biz, &aem)

		password := enterPassword("Export Mnemonic Words")
		if password == "" {
			executeWebview("msgbox", "-t", "Export Failed", "-c", "Please Enter account password", "-b", "Close")
			response.SetBaseResponse(errorcode.FAIL_CODE, "Please Enter account password")
			return response, nil
		}
		resp, err := service.ExportMnemonic(aem.Address, password)

		if err != nil {
			executeWebview("msgbox", "-t", "Export Failed!", "-c", err.Error(), "-b", "Close")
			response.SetBaseResponse(errorcode.FAIL_CODE, "The password is incorrect")
		} else {
			executeWebview("msgbox", "-t", "Export Successful, Write it down!", "-c", resp, "-b", "Close")
			response.SetBizResponse(true)
		}

		return response, nil
	}
}

// assets main
func MakeAccountListEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transport.Request)
		response := transport.Response{}
		response.SetBaseResponseSuccess()

		if ok, err := validator.ValidateBaseRequestParam(req.Base); !ok {
			response.SetBaseResponse(errorcode.InvalidBaseParameters, err.Error())
			return response, nil
		}
		response.SetBizResponse(service.AccountList())

		return response, nil
	}
}

type pk struct {
	PK string
}

// assets main
func MakeAccountDetailEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transport.Request)
		response := transport.Response{}
		response.SetBaseResponseSuccess()

		if ok, err := validator.ValidateBaseRequestParam(req.Base); !ok {
			response.SetBaseResponse(errorcode.InvalidBaseParameters, err.Error())
			return response, nil
		}
		pk := pk{}
		utils.Convert(req.Biz, &pk)
		ac := service.AccountDetail(pk.PK)
		response.SetBizResponse(ac)
		return response, nil
	}
}

func MakeAccountBalanceEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transport.Request)
		response := transport.Response{}
		response.SetBaseResponseSuccess()

		if ok, err := validator.ValidateBaseRequestParam(req.Base); !ok {
			response.SetBaseResponse(errorcode.InvalidBaseParameters, err.Error())
			return response, nil
		}
		pk := pk{}
		utils.Convert(req.Biz, &pk)
		balance,tickets := service.AccountBalance(pk.PK)
		//for key, v := range balance {
		//}
		resp := map[string]interface{}{}
		resp["balance"] = balance
		resp["tickets"] = tickets
		response.SetBizResponse(resp)
		return response, nil
	}
}

func MakeTxListEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transport.Request)
		response := transport.Response{}
		response.SetBaseResponseSuccess()

		if ok, err := validator.ValidateBaseRequestParam(req.Base); !ok {
			response.SetBaseResponse(errorcode.InvalidBaseParameters, err.Error())
			return response, nil
		}
		pk := pk{}
		utils.Convert(req.Biz, &pk)
		if records, err := service.TXList(pk.PK, req.Page); err != nil {
			response.SetBaseResponse(errorcode.FAIL_CODE, err.Error())
			return response, nil
		} else {
			response.SetPageResponse(uint8(len(records)), req.Page.PageSize, req.Page.PageNo, "")
			response.SetBizResponse(records)
		}
		return response, nil
	}
}

func MakeTxNumEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transport.Request)
		response := transport.Response{}
		response.SetBaseResponseSuccess()

		if ok, err := validator.ValidateBaseRequestParam(req.Base); !ok {
			response.SetBaseResponse(errorcode.InvalidBaseParameters, err.Error())
			return response, nil
		}
		pk := pk{}
		utils.Convert(req.Biz, &pk)
		response.SetBizResponse(service.TXNum(pk.PK))
		return response, nil
	}
}

func MakeTxSendEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transport.Request)
		response := transport.Response{}
		response.SetBaseResponseSuccess()

		if ok, err := validator.ValidateBaseRequestParam(req.Base); !ok {
			response.SetBaseResponse(errorcode.InvalidBaseParameters, err.Error())
			return response, nil
		}

		transferReq := transferReq{}
		utils.Convert(req.Biz, &transferReq)

		//password := enterPassword("Send Transfer")
		//if password == "" {
		//	response.SetBaseResponse(errorcode.FAIL_CODE, "Please Enter account password.")
		//	return response, nil
		//}

		hash, err := service.Transfer(transferReq, transferReq.Password)
		if err != nil {
			response.SetBaseResponse(errorcode.FAIL_CODE, err.Error())
			return response, nil
		}
		response.SetBizResponse(hash)
		return response, nil
	}
}

type transferReq struct {
	From     string
	To       string
	Currency string
	Amount   string
	Gas      string
	GasPrice string
	Password string
	Data     string
	AssetTktReq map[string]interface{} `json:"tkt"`
}

func MakeDataPathEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(transport.Request)
		response := transport.Response{}
		response.SetBaseResponseSuccess()

		if ok, err := validator.ValidateBaseRequestParam(req.Base); !ok {
			response.SetBaseResponse(errorcode.InvalidBaseParameters, err.Error())
			return response, nil
		}
		response.SetBizResponse(GetKeystorePath())
		return response, nil
	}
}

type decimalReq struct {
	Currency string
}

func MakeCurrencyDecimalEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(transport.Request)
		response := transport.Response{}
		response.SetBaseResponseSuccess()

		if ok, err := validator.ValidateBaseRequestParam(req.Base); !ok {
			response.SetBaseResponse(errorcode.InvalidBaseParameters, err.Error())
			return response, nil
		}
		decimalReq := decimalReq{}
		utils.Convert(req.Biz, &decimalReq)

		response.SetBizResponse(service.GetDecimal(decimalReq.Currency))
		return response, nil
	}
}

// stake pool
func MakeStakePoolEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(transport.Request)
		response := transport.Response{}
		response.SetBaseResponseSuccess()

		if ok, err := validator.ValidateBaseRequestParam(req.Base); !ok {
			response.SetBaseResponse(errorcode.InvalidBaseParameters, err.Error())
			return response, nil
		}

		sync := Sync{RpcHost: GetRpcHost(), Method: "stake_stakePools", Params: []interface{}{}}
		jsonResp, err := sync.Do()
		if err != nil {
			logex.Errorf("jsonRep err=[%s]", err.Error())
			return response, nil
		}
		if jsonResp.Result != nil {
			json.Unmarshal(*jsonResp.Result, &response.Biz)
		}
		return response, nil
	}
}

type registerReq struct {
	From     string
	Vote     string
	Password string
	FeeRate  string
	Type     string
	IdPkr    string
}

type closeStakeReq struct {
	From     string `validate:"required"`
	Password string
	IdPkr    string `validate:"required"`
}

//closeShare
func MakeCloseShareEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(transport.Request)
		response := transport.Response{}
		response.SetBaseResponseSuccess()

		if ok, err := validator.ValidateBaseRequestParam(req.Base); !ok {
			response.SetBaseResponse(errorcode.InvalidBaseParameters, err.Error())
			return response, nil
		}

		regO := closeStakeReq{}
		utils.Convert(req.Biz, &regO)

		//password := enterPassword("Close Stake Node")
		password := regO.Password
		if password == "" {
			response.SetBaseResponse(errorcode.FAIL_CODE, "Please Enter account password.")
			return response, nil
		}

		txHash, err := service.closeStake(regO.From, regO.IdPkr, password)
		if err != nil {
			response.SetBaseResponse(errorcode.FAIL_CODE, err.Error())
		} else {
			response.SetBizResponse(txHash)
		}
		return response, nil
	}
}

//registerShare
func MakeRegisterShareEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(transport.Request)
		response := transport.Response{}
		response.SetBaseResponseSuccess()

		if ok, err := validator.ValidateBaseRequestParam(req.Base); !ok {
			response.SetBaseResponse(errorcode.InvalidBaseParameters, err.Error())
			return response, nil
		}

		regO := registerReq{}
		utils.Convert(req.Biz, &regO)

		//password := enterPassword("Register or Update Stake Node")
		password := regO.Password
		if password == "" {
			response.SetBaseResponse(errorcode.FAIL_CODE, "Please Enter account password.")
			return response, nil
		}
		feeRate, err := strconv.ParseUint(regO.FeeRate, 10, 64)
		if err != nil {
			response.SetBaseResponse(errorcode.FAIL_CODE, "fee rate is two decimal places, eg: 25.55% ")
			return response, nil
		}
		if regO.Type == "" {
			txHash, err := service.registerStakePool(regO.From, regO.Vote, password, uint32(feeRate))
			if err != nil {
				response.SetBaseResponse(errorcode.FAIL_CODE, err.Error())
			} else {
				response.SetBizResponse(txHash)
			}
		} else if regO.Type == "modify" {
			txHash, err := service.modifyStakePool(regO.From, regO.Vote, password, regO.IdPkr, uint32(feeRate))
			if err != nil {
				response.SetBaseResponse(errorcode.FAIL_CODE, err.Error())
			} else {
				response.SetBizResponse(txHash)
			}
		} else {
			response.SetBaseResponse(errorcode.InvalidBaseParameters, "Type is invalid")
			return response, nil
		}

		return response, nil
	}
}

type buyShareReq struct {
	From     string
	Vote     string
	Password string
	Pool     string
	Amount   string
	GasPrice string
}

//buyShare
func MakeBuyShareEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(transport.Request)
		response := transport.Response{}
		response.SetBaseResponseSuccess()

		if ok, err := validator.ValidateBaseRequestParam(req.Base); !ok {
			response.SetBaseResponse(errorcode.InvalidBaseParameters, err.Error())
			return response, nil
		}

		regO := buyShareReq{}
		utils.Convert(req.Biz, &regO)
		//password := enterPassword("Buy Share")
		password := regO.Password
		if password == "" {
			response.SetBaseResponse(errorcode.FAIL_CODE, "Please Enter account password.")
			return response, nil
		}
		txHash, err := service.buyStake(regO.From, regO.Vote, password, regO.Pool, regO.Amount, regO.GasPrice)
		if err != nil {
			response.SetBaseResponse(errorcode.FAIL_CODE, err.Error())
		} else {
			response.SetBizResponse(txHash)
		}

		return response, nil
	}
}

//getShare
func MakeGetShareEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(transport.Request)
		response := transport.Response{}
		response.SetBaseResponseSuccess()

		if ok, err := validator.ValidateBaseRequestParam(req.Base); !ok {
			response.SetBaseResponse(errorcode.InvalidBaseParameters, err.Error())
			return response, nil
		}

		sync := Sync{RpcHost: GetRpcHost(), Method: "stake_getShare", Params: []interface{}{req.Biz}}
		jsonResp, err := sync.Do()
		if err != nil {
			return response, nil
		}
		if jsonResp.Result != nil {
			json.Unmarshal(*jsonResp.Result, &response.Biz)
		}
		return response, nil
	}
}

//myShare
func MakeGetMySharesEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(transport.Request)
		response := transport.Response{}
		response.SetBaseResponseSuccess()

		if ok, err := validator.ValidateBaseRequestParam(req.Base); !ok {
			response.SetBaseResponse(errorcode.InvalidBaseParameters, err.Error())
			return response, nil
		}

		sync := Sync{RpcHost: GetRpcHost(), Method: "stake_getShareByPkr", Params: []interface{}{req.Biz}}
		jsonResp, err := sync.Do()
		if err != nil {
			logex.Errorf("jsonRep err=[%s]", err.Error())
			return response, nil
		}
		if jsonResp.Result != nil {
			json.Unmarshal(*jsonResp.Result, &response.Biz)
		}
		return response, nil
	}
}

//getTransactionReceipt
func MakeGetTransactionReceiptEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(transport.Request)
		response := transport.Response{}
		response.SetBaseResponseSuccess()

		if ok, err := validator.ValidateBaseRequestParam(req.Base); !ok {
			response.SetBaseResponse(errorcode.InvalidBaseParameters, err.Error())
			return response, nil
		}

		sync := Sync{RpcHost: GetRpcHost(), Method: "sero_getTransactionReceipt", Params: []interface{}{req.Biz}}
		jsonResp, err := sync.Do()
		if err != nil {
			//logex.Errorf("jsonRep err=[%s]", err.Error())
			return response, nil
		}
		if jsonResp.Result != nil {
			json.Unmarshal(*jsonResp.Result, &response.Biz)
		}
		return response, nil
	}
}

//getTransactionReceipt
func MakeGetBlockNumberEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(transport.Request)
		response := transport.Response{}
		response.SetBaseResponseSuccess()

		if ok, err := validator.ValidateBaseRequestParam(req.Base); !ok {
			response.SetBaseResponse(errorcode.InvalidBaseParameters, err.Error())
			return response, nil
		}

		sync := Sync{RpcHost: GetRpcHost(), Method: "sero_blockNumber", Params: []interface{}{}}
		jsonResp, err := sync.Do()
		if err != nil {
			//logex.Errorf("jsonRep err=[%s]", err.Error())
			return response, nil
		}
		if jsonResp.Result != nil {
			json.Unmarshal(*jsonResp.Result, &response.Biz)
		}
		return response, nil
	}
}

func MakeChangeNetworkEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(transport.Request)
		response := transport.Response{}
		response.SetBaseResponseSuccess()

		if ok, err := validator.ValidateBaseRequestParam(req.Base); !ok {
			response.SetBaseResponse(errorcode.InvalidBaseParameters, err.Error())
			return response, nil
		}
		var hostReq = ""
		if req.Biz != nil {
			hostReq = req.Biz.(string)
		}
		response.SetBizResponse(service.getSetNetwork(hostReq))

		return response, nil
	}
}

func MakeOpenFileEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(transport.Request)
		response := transport.Response{}
		response.SetBaseResponseSuccess()

		if ok, err := validator.ValidateBaseRequestParam(req.Base); !ok {
			response.SetBaseResponse(errorcode.InvalidBaseParameters, err.Error())
			return response, nil
		}

		if GetOsType() == "mac" {
			c := "open " + app_home_path
			cmd := exec.Command("sh", "-c", c)
			_, err := cmd.Output()
			if err != nil {
				logex.Errorf("err:%s", err.Error())
			}
		} else if GetOsType() == "win" {
			cmd := exec.Command("explorer.exe", app_home_path)
			err := cmd.Run()
			if err != nil {
				logex.Errorf("open file err:%s", err.Error())
			}
		}

		return response, nil
	}
}

func MakeSetDappsEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(transport.Request)
		response := transport.Response{}
		response.SetBaseResponseSuccess()

		if ok, err := validator.ValidateBaseRequestParam(req.Base); !ok {
			response.SetBaseResponse(errorcode.InvalidBaseParameters, err.Error())
			return response, nil
		}
		dapp := Dapp{}
		utils.Convert(req.Biz, &dapp)
		rest, err := service.setDapps(dapp)
		if err != nil {
			response.SetBaseResponse(errorcode.FAIL_CODE, err.Error())
			return response, nil
		}
		response.SetBizResponse(rest)

		return response, nil
	}
}

type Dapp struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	Img       string `json:"img"`
	Title     string `json:"title"`
	Desc      string `json:"desc"`
	Author    string `json:"author"`
	Operation string `json:"operation"`
}
