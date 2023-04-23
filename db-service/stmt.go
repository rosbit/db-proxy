package dbsvc

import (
	"db-proxy/utils"
	"db-proxy/db"
	"net/http"
	"fmt"
)

// POST /close-stmt
// BODY: {"sid": "xx"}
func CloseStmt(svc DBService) {
	var params struct {
		StmtId string `json:"sid"`
	}
	if status, err := svc.ReadJSON(&params); err != nil {
		svc.SetErrorResp(status, err.Error())
		return
	}
	if len(params.StmtId) == 0 {
		svc.SetErrorResp(http.StatusBadRequest, "bad request")
		return
	}
	stmtId, idType, err := utils.DecodeId(params.StmtId)
	if err != nil {
		svc.SetErrorResp(http.StatusBadRequest, err.Error())
		return
	}
	if idType != utils.ID_STMT {
		svc.SetErrorResp(http.StatusBadRequest, "bad id type")
		return
	}
	stmt := db.Stmt(stmtId)
	if err := stmt.Close(); err != nil {
		svc.SetErrorResp(http.StatusInternalServerError, err.Error())
		return
	}

	svc.SetResult(nil)
}

// POST /exec
// BODY: {"sid": "xxx", "args":[xxxx, "adafasf"]}
func Exec(svc DBService) {
	stmt, args, status, err := getStmtAndArgs(svc)
	if err != nil {
		svc.SetErrorResp(status, err.Error())
		return
	}
	lastInsertId, rowsAffected, err := stmt.Exec(args...)
	if err != nil {
		svc.SetErrorResp(http.StatusInternalServerError, err.Error())
		return
	}
	svc.SetResult(map[string]interface{}{
		"lastInsertId": lastInsertId,
		"rowsAffected": rowsAffected,
	})
}

// POST /query
// BODY: {"sid": "xxx", "args":[xxx, "dafdas"]}
func Query(svc DBService) {
	stmt, args, status, err := getStmtAndArgs(svc)
	if err != nil {
		svc.SetErrorResp(status, err.Error())
		return
	}

	columns, _, it, err := stmt.Query(args...)
	if err != nil {
		svc.SetErrorResp(http.StatusInternalServerError, err.Error())
		return
	}

	svc.SetResult(map[string]interface{}{
		"columns": columns,
		// "columnTypes": columnTypes,
	}, it)
}

func getStmtAndArgs(svc DBService) (stmt db.Stmt, args []interface{}, status int, err error) {
	var params struct {
		StmtId string `json:"sid"`
		Args []interface{} `json:"args"`
	}
	if status, err = svc.ReadJSON(&params); err != nil {
		return
	}
	if len(params.StmtId) == 0 {
		status, err = http.StatusBadRequest, fmt.Errorf("bad request")
		return
	}
	stmtId, idType, e := utils.DecodeId(params.StmtId)
	if e != nil {
		status, err = http.StatusBadRequest, e
		return
	}
	if idType != utils.ID_STMT {
		status, err = http.StatusBadRequest, fmt.Errorf("bad id type")
		return
	}
	stmt = db.Stmt(stmtId)
	args = params.Args

	return
}
