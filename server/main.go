package main

import (
	"container/list"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

// code:web9-5
// 接続しているクライアントを保持
var conns list.List

func echo(ws *websocket.Conn) {
	fmt.Println(ws.Config().Origin)
	var err error
	// 新規のクライアントならconnに追加
	var conn *list.Element
	contain := false
	for c := conns.Front(); c != nil; c = c.Next() {
		if c.Value == ws {
			contain = true
		}
		fmt.Println("c value")
		fmt.Println(c.Value.(*websocket.Conn).Config().Origin)
	}
	if !contain {
		fmt.Println(ws.Config().Origin)
		fmt.Println("push back")
		conn = conns.PushBack(ws)
	}
	for {
		var reply string
		clientHost := ws.Config().Origin
		// 接続が途切れたらconnから削除
		if err = websocket.Message.Receive(ws, &reply); err != nil {
			conns.Remove(conn)
			ws.Close()
			fmt.Println("cant recieve", clientHost)
			break
		}
		// 送信するメッセージの生成
		fmt.Println("Recieve back from client:", reply)
		msg := "Recieved:" + reply + "from:" + clientHost.Host + "at:" + string(time.Now().Format("2006/1/2 15:04:05"))
		fmt.Println("Sending to client:", msg)
		// 保持している全てのクライアントに送信
		for c := conns.Front(); c != nil; c = c.Next() {
			fmt.Println("push from")
			fmt.Println(ws.Config().Origin)
			fmt.Println("push to")
			fmt.Println(c.Value.(*websocket.Conn).Config().Origin)
			if err = websocket.Message.Send(c.Value.(*websocket.Conn), msg); err != nil {
				fmt.Println("cant send.")
				break
			}
			reply = ""
		}
	}
}

func main() {
	http.Handle("/", websocket.Handler(echo))

	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("Listen and Serve:", err)
	}
}
