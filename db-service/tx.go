package dbsvc

import (
	"db-proxy/utils"
	"db-proxy/db"
	"net/http"
	"fmt"
)

// POST /rollback
// BODY: {"tid": "xxx"}
func Rollback(svc DBService) {
	tx, status, err := getTx(svc)
	if err != nil {
		svc.SetErrorResp(status, err.Error())
		return
	}
	endTx(svc, tx.Rollback)
}

// POST /commit
// BODY: {"tid": "xxx"}
func Commit(svc DBService) {
	tx, status, err := getTx(svc)
	if err != nil {
		svc.SetErrorResp(status, err.Error())
		return
	}
	endTx(svc, tx.Commit)
}

func endTx(svc DBService, end func()error) {
	if err := end(); err != nil {
		svc.SetErrorResp(http.StatusInternalServerError, err.Error())
		return
	}
	svc.SetResult(nil)
}

func getTx(svc DBService) (tx db.Tx, status int, err error) {
	var params struct {
		TxId string `json:"tid"`
	}
	if status, err = svc.ReadJSON(&params); err != nil {
		return
	}
	if len(params.TxId) == 0 {
		status, err = http.StatusBadRequest, fmt.Errorf("bad request")
		return
	}
	txId, idType, e := utils.DecodeId(params.TxId)
	if e != nil {
		status, err = http.StatusBadRequest, e
		return
	}
	if idType != utils.ID_TX {
		status, err = http.StatusBadRequest, fmt.Errorf("bad id type")
		return
	}
	tx = db.Tx(txId)
	return
}
