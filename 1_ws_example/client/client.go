package main

import (
	"log"

	"github.com/gorilla/websocket"
)

func main() {
	url := "ws://localhost:8080"

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("连接失败:", err)
	}
	defer conn.Close()
	log.Println("成功连接到", url)

	message := "Hello from Go client!"
	err = conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Fatal("发送消息失败:", err)
	}
	log.Println("已发送消息:", message)

	// 接收服务器返回的消息
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Fatal("读取消息失败:", err)
	}
	log.Println("接收到服务器消息:", string(msg))
}
