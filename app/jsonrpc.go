package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sero-cash/go-sero/light-wallet/common/logex"
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
	Code int `json:"code"`
	Message string `json:"message"`
}

func (sync Sync) Do() (*JSONRpcResp, error) {
	client := &http.Client{
		Timeout: 900 * time.Second,
	}
	//logex.Info("sync.Params=", sync.Params)
	jsonReq := map[string]interface{}{"jsonrpc": "2.0", "method": sync.Method, "params": sync.Params, "id": 0}
	data, err := json.Marshal(jsonReq)
	if err != nil {
		logex.Error(err.Error())
		return nil, err
	}
	logex.Info(string(data))

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
