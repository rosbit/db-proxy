package rest

import (
	"db-proxy/db-service"
	"net"
	"fmt"
	"net/http"
)

func CreateDBProxy(reqAction string) func(c *Context) {
	return func(c *Context) {
		dbService := createDBService(c)
		action, ok := dbsvc.GetAction(reqAction)
		if !ok {
			dbService.SetErrorResp(http.StatusNotFound, "action not found")
			return
		}
		action(dbService)
	}
}

func TryListen(host string, port int) (l net.Listener, err error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	if l, err = net.Listen("tcp", addr); err != nil {
		return
	}
	fmt.Printf("I am listening at %s\n", addr)
	return
}
