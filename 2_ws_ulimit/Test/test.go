package main

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	const connections = 1000
	var wg sync.WaitGroup
	wg.Add(connections)

	for i := range connections {
		go func(i int) {
			defer wg.Done()
			// 连接 WebSocket 服务
			conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8000/", nil)
			if err != nil {
				log.Printf("连接失败: %v", err)
				return
			}
			defer conn.Close()

			// 发送一条消息
			if err := conn.WriteMessage(websocket.TextMessage, []byte("hello")); err != nil {
				log.Printf("发送消息失败: %v", err)
			}

			// 保持连接 30 秒钟
			time.Sleep(30 * time.Second)
		}(i)
	}

	wg.Wait()
	log.Println("所有连接已关闭")
}
