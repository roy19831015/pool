package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/roy19831015/pool/pool"
	"net/http"
	"os"
	"strconv"
	"time"
)

type HttpHandler struct {
	http.Handler
}

func (httpHandler HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("程序进入监听>>")
	var upgrader = websocket.Upgrader{
		//解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("upgrade error:", err)
		return
	}
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(6 * time.Second))
	conn.SetPongHandler(func(string) error {
		fmt.Println("接收心跳响应<<")
		conn.SetReadDeadline(time.Now().Add(6 * time.Second))
		return nil
	})
	ticker := time.NewTicker(2 * time.Second)
	go func() {
		for {
			//从定时器中获取数据
			_ = <-ticker.C
			fmt.Println("发送心跳包>>")
			conn.WriteMessage(websocket.PingMessage, []byte{})
		}
	}()
	defer ticker.Stop()
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("接收异常：", err, "消息类型：", messageType)
			break
		}
		fmt.Println(fmt.Sprintf("接收消息内容 >>%s", message))
		respMessage := fmt.Sprintf("I am Server, %s.", time.Now().Format("2006-01-02 15:04:05"))
		err = conn.WriteMessage(messageType, []byte(respMessage))
		if err != nil {
			fmt.Println("发送异常：", err)
			break
		}
	}
}

func main() {
	go func() {
		var httpHandler HttpHandler
		http.Handle("/", httpHandler)
		if err := http.ListenAndServe(":3456", nil); err != nil {
			fmt.Println("程序退出")
			os.Exit(1)
		}
	}()
	time.After(time.Second * 5)
	p := pool.CommonPool[*websocket.Conn]{}
	err := p.Init(10, func() (*websocket.Conn, error) {
		conn, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:3456/", nil)
		if err != nil {
			return nil, err
		}
		conn.SetPingHandler(func(appData string) error {
			err := conn.WriteMessage(websocket.PongMessage, []byte("Heart Beat"))
			if err != nil {
				err := conn.Close()
				if err != nil {
					return nil
				}
			}
			return nil
		})
		conn.SetCloseHandler(func(code int, text string) error {
			err := conn.Close()
			if err != nil {
				return err
			}
			return nil
		})
		return conn, nil
	})
	if err != nil {
		return
	}

	for i := 0; i < 100; i++ {
		go func(i int) {
			for {
				conn, err := p.Get(time.Second)
				if err != nil {
					return
				}
				err = conn.WriteMessage(websocket.TextMessage, []byte("client demo message"))
				if err != nil {
					err := conn.Close()
					if err != nil {
						return
					}
				}
				_, data, err := conn.ReadMessage()
				if err != nil {
					return
				}
				p.Back(conn)
				fmt.Println("message receive from server:", string(data), " No.", strconv.Itoa(i))
			}
		}(i)
	}
	ch := make(chan int)
	<-ch
}
