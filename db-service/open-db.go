package dbsvc

import (
	"db-proxy/utils"
	"db-proxy/db"
	"net/http"
)

// POST /open-db
// BODY: {"dsn": "<driver-name>:<dsn-string>"}
func OpenDB(svc DBService) {
	var params struct {
		DSN string `json:"dsn"`
	}
	if status, err := svc.ReadJSON(&params); err != nil {
		svc.SetErrorResp(status, err.Error())
		return
	}
	if len(params.DSN) == 0 {
		svc.SetErrorResp(http.StatusBadRequest, "bad request")
		return
	}
	dsn, err := db.GetDSNByName(params.DSN)
	if err != nil {
		svc.SetErrorResp(http.StatusInternalServerError, err.Error())
		return
	}

	if _, err = db.Open(dsn); err != nil {
		svc.SetErrorResp(http.StatusInternalServerError, err.Error())
		return
	}

	svc.SetResult(map[string]interface{}{
		"db-id": utils.EncodeStr(dsn),
	})
}
