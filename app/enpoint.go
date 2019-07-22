package app

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"github.com/sero-cash/go-sero/pullup/common/errorcode"
	"github.com/sero-cash/go-sero/pullup/common/logex"
	"github.com/sero-cash/go-sero/pullup/common/transport"
	"github.com/sero-cash/go-sero/pullup/common/utils"
	"github.com/sero-cash/go-sero/pullup/common/validator"
	"os/exec"
	"strconv"
)

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

		resp, err := service.NewAccountWithMnemonic(accountCreateReq.Passphrase)
		if err != nil {
			response.SetBaseResponse(errorcode.FAIL_CODE, err.Error())
		} else {
			response.SetBizResponse(resp)
		}
		return response, nil
	}
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

		resp, err := service.ImportAccountFromRawKey(aipq.PrivateKey, aipq.Passphrase)
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

		resp, err := service.ExportMnemonic(aem.Address, aem.Passphrase)
		if err != nil {
			response.SetBaseResponse(errorcode.FAIL_CODE, err.Error())
		} else {
			response.SetBizResponse(resp)
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
		balance := service.AccountBalance(pk.PK)
		//for key, v := range balance {
		//}
		response.SetBizResponse(balance)

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

		if transferReq.Password == "" {
			response.SetBaseResponse(errorcode.FAIL_CODE, "password can not be nil ")
			return response, nil
		}

		hash, err := service.Transfer(transferReq.From, transferReq.To, transferReq.Currency, transferReq.Amount, transferReq.GasPrice, transferReq.Password)
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
	GasPrice string
	Password string
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

		sync := Sync{RpcHost: host, Method: "stake_stakePools", Params: []interface{}{}}
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

		feeRate, err := strconv.ParseUint(regO.FeeRate, 10, 64)
		if err != nil {
			response.SetBaseResponse(errorcode.FAIL_CODE, "fee rate is two decimal places, eg: 25.55% ")
			return response, nil
		}

		txHash, err := service.registerStakePool(regO.From, regO.Vote, regO.Password, uint32(feeRate))
		if err != nil {
			response.SetBaseResponse(errorcode.FAIL_CODE, err.Error())
		} else {
			response.SetBizResponse(txHash)
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

		txHash, err := service.buyStake(regO.From, regO.Vote, regO.Password, regO.Pool, regO.Amount, regO.GasPrice)
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

		sync := Sync{RpcHost: host, Method: "stake_getShare", Params: []interface{}{req.Biz}}
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

		sync := Sync{RpcHost: host, Method: "stake_getShareByPkr", Params: []interface{}{req.Biz}}
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

		sync := Sync{RpcHost: host, Method: "sero_getTransactionReceipt", Params: []interface{}{req.Biz}}
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

		sync := Sync{RpcHost: host, Method: "sero_blockNumber", Params: []interface{}{}}
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
		host = service.getSetNetwork(hostReq)
		response.SetBizResponse(host)

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
			c := "open " + GetDataPath()
			cmd := exec.Command("sh", "-c", c)
			_, err := cmd.Output()
			if err != nil {
				logex.Errorf("err:", err.Error())
			}
		} else if GetOsType() == "win" {
			cmd := exec.Command("explorer.exe", GetDataPath())
			err :=cmd.Run()
			if err != nil {
				logex.Errorf("open file err:", err.Error())
			}
		}

		return response, nil
	}
}
