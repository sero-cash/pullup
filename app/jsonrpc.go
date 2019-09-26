package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sero-cash/go-sero/pullup/common/logex"
	"github.com/sero-cash/go-sero/rlp"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// === jsonrpc post

type Sync struct {
	RpcHost string
	Method  string
	Params  interface{}
}

type ErrorMsg struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (sync Sync) Do() (*JSONRpcResp, error) {
	client := &http.Client{
		Timeout: 900 * time.Second,
	}

	jsonReq := map[string]interface{}{"jsonrpc": "2.0", "method": sync.Method, "params": sync.Params, "id": 0}
	data, err := json.Marshal(jsonReq)
	if err != nil {
		logex.Error(err.Error())
		return nil, err
	}

	req, err := http.NewRequest("POST", sync.RpcHost, bytes.NewBuffer(data))
	if err != nil {
		logex.Error(err.Error())
		return nil, err
	}
	req.Header.Set("Content-Length", (string)(len(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		logex.Error(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	var rpcResp *JSONRpcResp
	err = json.NewDecoder(resp.Body).Decode(&rpcResp)
	if err != nil {
		logex.Error(err.Error())
		return nil, err
	}
	if rpcResp.Error != nil {
		logex.Error(rpcResp.Error)
		return nil, fmt.Errorf(rpcResp.Error["message"].(string))
	}
	return rpcResp, err
}

type JSONRpcResp struct {
	Id      *json.RawMessage       `json:"id"`
	Jsonrpc string                 `json:"jsonrpc"`
	Result  *json.RawMessage       `json:"result"`
	Error   map[string]interface{} `json:"error"`
}

type JSONRpcReturn struct {
	Id      json.RawMessage `json:"id"`
	Jsonrpc string          `json:"jsonrpc"`
	Result  interface{}     `json:"result"`
	Error   interface{}     `json:"error"`
}

type JSONRpcReq struct {
	Id      *json.RawMessage  `json:"id"`
	Jsonrpc string            `json:"jsonrpc"`
	Method  PULLUP_RPC_METHOD `json:"method"`
	Params  json.RawMessage   `json:"params"`
}

// === pullup rpc handler
func HandlePullupRpc(req JSONRpcReq, service Service) (resp JSONRpcReturn) {
	resp.Id = *req.Id
	resp.Jsonrpc = "2.0"
	switch req.Method {
	case RPC_METHOD_DEPLOY_CONTRACT:
		reqTx := sentContractTxReq{}
		err := json.Unmarshal(req.Params[:], &reqTx)
		if err != nil {
			resp.Error = Reply{Code: -1, Message: err.Error()}
		}
		password := enterPassword("Deploy Contract")
		if password == "" {
			resp.Error = Reply{Code: -1, Message: "Please enter account password"}
			return resp
		}
		txNo, err := currentLight.DeployContractTx(reqTx.ContractTxReq, password)
		if err != nil {
			resp.Error = Reply{Code: -1, Message: err.Error()}
			return resp
		}
		resp.Result = txNo
		break
	case RPC_METHOD_EXECUTE_CONTRACT:
		reqTx := sentContractTxReq{}
		err := json.Unmarshal(req.Params[:], &reqTx)
		if err != nil {
			resp.Error = Reply{Code: -1, Message: err.Error()}
		}
		password := enterPassword("Execute Contract")
		if password == "" {
			resp.Error = Reply{Code: -1, Message: "Please enter account password"}
			return resp
		}
		txHash, err := currentLight.ExecuteContractTx(reqTx.ContractTxReq, password)
		if err != nil {
			resp.Error = Reply{Code: -1, Message: err.Error()}
			return resp
		}
		resp.Result = txHash
		break
	case EXECUTE_CONTRACT:
		reqTx := []ContractTxReq{}

		fmt.Println(string(req.Params[:]))
		err := json.Unmarshal(req.Params, &reqTx)
		if err != nil {
			resp.Error = Reply{Code: -1, Message: err.Error()}
		}
		password := enterPassword("Execute Contract")
		if password == "" {
			resp.Error = Reply{Code: -1, Message: "Please enter account password"}
			return resp
		}
		txHash, err := currentLight.ExecuteContractTx(reqTx[0], password)
		if err != nil {
			resp.Error = Reply{Code: -1, Message: err.Error()}
			return resp
		}
		resp.Result = txHash
		break
	case RPC_METHOD_GET_TOENS:
		tokens, err := currentLight.getTokens()
		if err != nil {
			resp.Error = Reply{Code: -1, Message: err.Error()}
			return resp
		}
		resp.Result = tokens
		break
	case RPC_METHOD_WATCH_TOENS:
		token := TokenReq{}
		err := json.Unmarshal(req.Params[:], &token)
		if err != nil {
			resp.Error = Reply{Code: -1, Message: err.Error()}
			return resp
		}
		data, err := rlp.EncodeToBytes(token)
		if err != nil {
			resp.Error = Reply{Code: -1, Message: err.Error()}
			return resp
		}
		currentLight.db.Put(append(tokenPrefix, []byte(token.ContractAddress)[:]...), data[:])
		resp.Result = true
		break
	case ACCOUNT_LIST:
		resp.Result = service.AccountList()
		break
	case ACCOUNT_DETAIL:
		params := []string{}
		json.Unmarshal(req.Params[:], &params)
		//fmt.Println("string(data[:])== ", string(params[0]))
		resp.Result = service.AccountDetail(string(params[0]))
		break
	}

	return resp
}

//
type Reply struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
}

type sentContractTxReq struct {
	Password      string        `json:"password"`
	ContractTxReq ContractTxReq `json:"contract_tx_req"`
}

type PULLUP_RPC_METHOD string

var (
	RPC_METHOD_DEPLOY_CONTRACT  PULLUP_RPC_METHOD = "deploy_contract"
	RPC_METHOD_EXECUTE_CONTRACT PULLUP_RPC_METHOD = "execute_contract"
	RPC_METHOD_GET_TOENS        PULLUP_RPC_METHOD = "get_tokens"
	RPC_METHOD_WATCH_TOENS      PULLUP_RPC_METHOD = "watch_tokens"

	EXECUTE_CONTRACT PULLUP_RPC_METHOD = "pullup_execute_contract"
	ACCOUNT_LIST     PULLUP_RPC_METHOD = "pullup_account_list"
	ACCOUNT_DETAIL   PULLUP_RPC_METHOD = "pullup_account_detail"
)

func DoRequest(url string) (body []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("response err:", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	log.Println("response json:", string(body))
	return body, err
}