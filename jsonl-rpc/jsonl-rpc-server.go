package rpc

import (
	"github.com/rosbit/jsonl-rpc"
	"db-proxy/db-service"
	r "net/rpc"
	"runtime"
	"time"
	"net"
	"os"
	"fmt"
	"log"
	"net/http"
	"encoding/json"
)

type DBProxyRequest struct {
	Action string          `json:"action"`
	Args   json.RawMessage `json:"args"`
}

type JSONLServerResponse = jsonlrpc.JSONLServerResponse

type DBProxy int
func (p *DBProxy) Do(req *DBProxyRequest, res *JSONLServerResponse) (err error) {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		log.Printf("[jsonlrpc-db-proxy] %10s => %d %10s %s\n", req.Action, res.Code, duration.String(), res.Msg)
	}()

	dbService := createDBService(req.Args, res)
	action, ok := dbsvc.GetAction(req.Action)
	if !ok {
		dbService.SetErrorResp(http.StatusNotFound, "action not found")
		return
	}

	action(dbService)
	return
}

func TryListen(host string, port int) (l net.Listener, err error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	if l, err = net.Listen("tcp", addr); err != nil {
		return
	}
	fmt.Printf("jsonl-rpc-server is listening at %s\n", addr)
	return
}

func init() {
	dbProxy := new(DBProxy)
	r.Register(dbProxy)
}

func StartJSONLRpcServer(l net.Listener) {
	defer l.Close()

	for {
		conn, e := l.Accept()
		if e != nil {
			break
		}

		go func() {
			// fmt.Fprintf(os.Stderr, "connection accepted\n")
			defer func() {
				// fmt.Fprintf(os.Stderr, "connection closed\n")
				if err := recover(); err != nil {
					stack := make([]byte, 8*1024)
					size := runtime.Stack(stack, true)
					fmt.Fprintf(os.Stderr, "%s\n", stack[:size])
				}
			}()

			jsonlrpc.ServeConn(conn)
		}()
	}
}
