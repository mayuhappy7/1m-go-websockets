package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"sync/atomic"
	"syscall"

	"github.com/gorilla/websocket"
)

var count int64

func ws(w http.ResponseWriter, r *http.Request) {
	// Upgrade connection
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	n := atomic.AddInt64(&count, 1)
	if n%100 == 0 {
		log.Printf("Total number of connections: %v", n)
	}
	defer func() {
		n := atomic.AddInt64(&count, -1)
		if n%100 == 0 {
			log.Printf("Total number of connections: %v", n)
		}
		conn.Close()
	}()

	// Read messages from socket
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			return
		}
		// log.Printf("msg: %s", string(msg))
		// // 回写消息给客户端，实现 echo 功能
		// if err := conn.WriteMessage(messageType, msg); err != nil {
		// 	return
		// }
	}
}

func main() {

	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

	// Enable pprof hooks
	go func() {
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatalf("Pprof failed: %v", err)
		}
	}()

	http.HandleFunc("/", ws)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
