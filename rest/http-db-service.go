package rest

import (
	"github.com/rosbit/mgin"
	"net/http"
	"encoding/json"
)

type Context = mgin.Context

type httpDBService struct {
	c *Context
}
func createDBService(c *Context) *httpDBService {
	return &httpDBService{
		c: c,
	}
}

// implementation of db-service/DBService
func (svc *httpDBService) SetErrorResp(status int, msg string) {
	svc.c.Error(status, msg)
}
func (svc *httpDBService) SetResult(result map[string]interface{}, rows ...<-chan interface{}) {
	svc.c.JSON(http.StatusOK, map[string]interface{} {
		"code": http.StatusOK,
		"msg": "OK",
		"result": result,
	})
	if len(rows) == 0 || rows[0] == nil {
		return
	}

	out := json.NewEncoder(svc.c.Response())
	out.SetEscapeHTML(false)
	for row := range rows[0] {
		out.Encode(row)
	}
}
func (svc *httpDBService) ReadJSON(params interface{}) (status int, err error) {
	return svc.c.ReadJSON(params)
}
