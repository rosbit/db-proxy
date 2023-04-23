package rpc

import (
	"encoding/json"
	"net/http"
	// "fmt"
)

type jsonlDBService struct {
	args json.RawMessage
	res *JSONLServerResponse
}
func createDBService(args json.RawMessage, res *JSONLServerResponse) *jsonlDBService {
	return &jsonlDBService{
		args: args,
		res: res,
	}
}

// implementation of db-service/DBService
func (svc *jsonlDBService) SetErrorResp(status int, msg string) {
	res := svc.res

	res.Code = status
	res.Msg  = msg
	res.Result = nil
	res.HasJSONLs = false
	res.JSONLs = nil
}
func (svc *jsonlDBService) SetResult(result map[string]interface{}, rows ...<-chan interface{}) {
	res := svc.res

	res.Code = http.StatusOK
	res.Msg = "OK"
	res.Result = result
	if len(rows) > 0 && rows[0] != nil {
		res.HasJSONLs = true
		res.JSONLs = rows[0]
	} else {
		res.HasJSONLs = false
		res.JSONLs = nil
	}
}
func (svc *jsonlDBService) ReadJSON(params interface{}) (status int, err error) {
	// fmt.Printf("args: %s\n", string(svc.args))
	if err = json.Unmarshal(svc.args, params); err != nil {
		status = http.StatusBadRequest
	} else {
		status = http.StatusOK
	}
	return
}
