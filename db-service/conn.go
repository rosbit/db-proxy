package dbsvc

import (
	"db-proxy/utils"
	"db-proxy/db"
	"net/http"
	"fmt"
)

type iTxParam interface{
	getDBId() string
	getTxId() string
}
type baseParam struct {
	DbId string `json:"did"`
	TxId string `json:"tid"`
}
func (p *baseParam) getDBId() string {
	return p.DbId
}
func (p *baseParam) getTxId() string {
	return p.TxId
}

// POST /prepare
// BODY: {"did": "xx", "tid": "<optional>", "query":"xxx"}
func Prepare(svc DBService) {
	var params struct {
		baseParam
		Query  string `json:"query"`
	}
	dbInst, tx, status, err := getTxParams(svc, &params)
	if err != nil {
		svc.SetErrorResp(status, err.Error())
		return
	}
	if len(params.Query) == 0 {
		svc.SetErrorResp(http.StatusBadRequest, "bad request")
		return
	}

	var stmtId db.Stmt
	if len(params.TxId) == 0 {
		// DB.Parepare
		stmtId, err = dbInst.Prepare(params.Query)
	} else {
		stmtId, err = tx.Prepare(params.Query)
	}

	if err != nil {
		svc.SetErrorResp(http.StatusInternalServerError, err.Error())
		return
	}

	svc.SetResult(map[string]interface{}{
		"stmt-id": utils.EncodeId(utils.ID_STMT, uint32(stmtId)),
	})
	return
}

// POST /ping
// BODY: {"did": xxx, "tid": "<>"}
func Ping(svc DBService) {
	var params baseParam
	dbInst, tx, status, err := getTxParams(svc, &params)
	if err != nil {
		svc.SetErrorResp(status, err.Error())
		return
	}

	if len(params.TxId) == 0 {
		err = dbInst.Ping()
	} else {
		err = tx.Ping()
	}
	if err != nil {
		svc.SetErrorResp(http.StatusGone, err.Error())
		return
	}

	svc.SetResult(nil)
}

func getTxParams(svc DBService, params iTxParam) (dbInst db.DB, tx db.Tx, status int, err error) {
	if status, err = svc.ReadJSON(&params); err != nil {
		return
	}
	dbId := params.getDBId()
	if len(dbId) == 0 {
		status, err = http.StatusBadRequest, fmt.Errorf("bad request")
		return
	}
	dsn, e:= utils.DecodeStr(dbId)
	if e != nil {
		status, err = http.StatusBadRequest, e
		return
	}

	dId, ok := utils.GetIdByRef(dsn)
	if !ok {
		// reopen dsn
		dbInst, err = db.Open(dsn)
		if err != nil {
			status = http.StatusBadRequest
			return
		}
	} else {
		dbInst = db.DB(dId)
	}

	tId := params.getTxId()
	if len(tId) == 0 {
		return
	}
	txId, idType, e := utils.DecodeId(tId)
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

