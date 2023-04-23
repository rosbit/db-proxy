/**
 * REST API router
 * Rosbit Xu
 */
package main

import (
	"github.com/rosbit/mgin"
	"db-proxy/db-service"
	"db-proxy/jsonl-rpc"
	"db-proxy/utils"
	"db-proxy/conf"
	"db-proxy/rest"
	"db-proxy/db"
	"net/http"
	"fmt"
)

func StartService() error {
	serviceConf := conf.ServiceConf
	if err := utils.InitIdCodec(serviceConf.Base32Chars); err != nil {
		return err
	}

	api := mgin.NewMgin(mgin.WithLogger("db-proxy"))
	for action := range dbsvc.GetActionNames() {
		api.POST(fmt.Sprintf("/%s", action), rest.CreateDBProxy(action))
	}

	// health check
	api.GET("/health", func(c *mgin.Context) {
		c.String(http.StatusOK, "OK\n")
	})
	api.Get("/websocket", rpc.WebsocketRpcHandler)

	db.InitQ(serviceConf.QLen)

	restListenr, err := rest.TryListen(serviceConf.ListenHost, serviceConf.HttpListenPort)
	if err != nil {
		return err
	}
	rpcListener, err := rpc.TryListen(serviceConf.ListenHost, serviceConf.JSONLRpcListenPort)
	if err != nil {
		return err
	}

	done := make(chan struct{})
	go func() {
		http.Serve(restListenr, api)
	}()
	go rpc.StartJSONLRpcServer(rpcListener)
	<-done

	return nil
}

