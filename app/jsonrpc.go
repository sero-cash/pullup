package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sero-cash/go-sero/pullup/common/logex"
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
	logex.Info("sync.Params=", jsonReq)
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
	Id     *json.RawMessage       `json:"id"`
	Result *json.RawMessage       `json:"result"`
	Error  map[string]interface{} `json:"error"`
}

type JSONRpcReturn struct {
	Id     json.RawMessage `json:"id"`
	Result interface{}     `json:"result"`
	Error  interface{}     `json:"error"`
}

type JSONRpcReq struct {
	Id     *json.RawMessage  `json:"id"`
	Method PULLUP_RPC_METHOD `json:"method"`
	Params json.RawMessage   `json:"params"`
}

// === pullup rpc handler
func HandlePullupRpc(req JSONRpcReq) (resp JSONRpcReturn) {
	resp.Id = *req.Id
	switch req.Method {
	case RPC_METHOD_GenTxNo:
		ctq := ContractTxReq{}
		err := json.Unmarshal(req.Params[:], &ctq)
		if err != nil {
			resp.Error = Reply{Code: -1, Message: err.Error()}
			return resp
		}
		txNo, err := currentLight.GenTxNo(ctq)
		if err != nil {
			resp.Error = Reply{Code: -1, Message: err.Error()}
			return resp
		}
		resp.Result = txNo
		break
	case RPC_METHOD_GET_TX:
		reqT := getTxReq{}
		err := json.Unmarshal(req.Params[:], &reqT)
		if err != nil {
			resp.Error = Reply{Code: -1, Message: err.Error()}
			return resp
		}
		ctq, err := currentLight.GetPreSendTx(reqT.TxNo)
		if err != nil {
			resp.Error = Reply{Code: -1, Message: err.Error()}
			return resp
		}
		resp.Result = ctq
		break
	case RPC_METHOD_SEND_TX:
		reqTx := sentTxReq{}
		err := json.Unmarshal(req.Params[:], &reqTx)
		if err != nil {
			resp.Error = Reply{Code: -1, Message: err.Error()}
		}
		txHash, err := currentLight.DeployContractTx(reqTx.TxNo, reqTx.Password)
		if err != nil {
			resp.Error = Reply{Code: -1, Message: err.Error()}
			return resp
		}
		resp.Result = txHash
		break
	case RPC_METHOD_EXECUTE_CONTRACT:
		reqTx := sentTxReq{}
		err := json.Unmarshal(req.Params[:], &reqTx)
		if err != nil {
			resp.Error = Reply{Code: -1, Message: err.Error()}
		}
		txHash, err := currentLight.ExecuteContractTx(reqTx.TxNo, reqTx.Password)
		if err != nil {
			resp.Error = Reply{Code: -1, Message: err.Error()}
			return resp
		}
		resp.Result = txHash
		break
	}
	return resp
}

//
type Reply struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
}

type sentTxReq struct {
	Password string `json:"password"`
	TxNo     string `json:"tx_no"`
}

type executeContractReq struct {
	Password   string `json:"password"`
	TxNo       string `json:"tx_no"`
	MethodName string `json:"method_name"`
}

type getTxReq struct {
	TxNo string `json:"tx_no"`
}

type PULLUP_RPC_METHOD string

var (
	RPC_METHOD_GenTxNo          PULLUP_RPC_METHOD = "gen_tx_no"
	RPC_METHOD_GET_TX           PULLUP_RPC_METHOD = "get_tx"
	RPC_METHOD_SEND_TX          PULLUP_RPC_METHOD = "send_tx"
	RPC_METHOD_EXECUTE_CONTRACT PULLUP_RPC_METHOD = "execute_contract"
)
