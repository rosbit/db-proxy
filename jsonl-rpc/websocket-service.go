package rpc

import (
	"github.com/gorilla/websocket"
	"github.com/rosbit/jsonl-rpc"
	"net/http"
	"log"
	"runtime"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool{
		return true
	},
} // use default options

func WebsocketRpcHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade: %v\n", err)
		return
	}

	defer c.Close()

	// fmt.Fprintf(os.Stderr, "connection accepted\n")
	defer func() {
		// fmt.Fprintf(os.Stderr, "connection closed\n")
		if err := recover(); err != nil {
			stack := make([]byte, 8*1024)
			size := runtime.Stack(stack, true)
			log.Printf("%s\n", stack[:size])
		}
	}()

	conn := c.UnderlyingConn()
	jsonlrpc.ServeConn(conn)
}
