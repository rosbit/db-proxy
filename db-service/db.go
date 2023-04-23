package dbsvc

import (
	"db-proxy/utils"
	"db-proxy/db"
	"net/http"
	"fmt"
)

// POST /close-db
// BODY: {"did": "xxx"}
func CloseDB(svc DBService) {
	dbInst, status, err := getDB(svc, true)
	if err != nil {
		svc.SetErrorResp(status, err.Error())
		return
	}
	dbInst.Close()

	svc.SetResult(nil)
}

func getDB(svc DBService, dontReopen ...bool) (dbInst db.DB, status int, err error) {
	var params struct {
		DbId string `json:"did"`
	}
	if status, err = svc.ReadJSON(&params); err != nil {
		return
	}
	if len(params.DbId) == 0 {
		status, err = http.StatusBadRequest, fmt.Errorf("bad request")
		return
	}
	/*
	dbId, idType, e := utils.DecodeId(params.DbId)
	if e != nil {
		status, err = http.StatusBadRequest, e
		return
	}
	if idType != utils.ID_CONN {
		status, err = http.StatusBadRequest, fmt.Errorf("bad id type")
		return
	}*/
	dsn, e := utils.DecodeStr(params.DbId)
	if e != nil {
		status, err = http.StatusBadRequest, e
		return
	}
	dbId, ok := utils.GetIdByRef(dsn)
	if !ok {
		if len(dontReopen) > 0 && dontReopen[0] {
			status, err = http.StatusBadRequest, fmt.Errorf("bad dsn")
			return
		}
		dbInst, err = db.Open(dsn)
		if err != nil {
			status = http.StatusBadRequest
			return
		}
	} else {
		dbInst = db.DB(dbId)
	}
	return
}
