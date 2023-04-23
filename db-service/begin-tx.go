package dbsvc

import (
	"db-proxy/utils"
	"db-proxy/db"
	"net/http"
	"database/sql"
)

// POST /begin-tx
// BODY: {"did": "xx", "opts": {}}
func BeginTx(svc DBService) {
	var params struct {
		DbId string `json:"did"`
		Opts sql.TxOptions `json:"opts"`
	}
	if status, err := svc.ReadJSON(&params); err != nil {
		svc.SetErrorResp(status, err.Error())
		return
	}
	if len(params.DbId) == 0 {
		svc.SetErrorResp(http.StatusBadRequest, "bad request")
		return
	}
	dsn, err := utils.DecodeStr(params.DbId)
	if err != nil {
		svc.SetErrorResp(http.StatusBadRequest, err.Error())
		return
	}

	var dbInst db.DB
	dbId, ok := utils.GetIdByRef(dsn)
	if !ok {
		// reopen dsn
		dbInst, err = db.Open(dsn)
		if err != nil {
			svc.SetErrorResp(http.StatusBadRequest, err.Error())
			return
		}
	} else {
		dbInst = db.DB(dbId)
	}

	txId, err := dbInst.BeginTx(&params.Opts)
	if err != nil {
		svc.SetErrorResp(http.StatusInternalServerError, err.Error())
		return
	}

	svc.SetResult(map[string]interface{}{
		"tx-id": utils.EncodeId(utils.ID_TX, uint32(txId)),
	})
}
