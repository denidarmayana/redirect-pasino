package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

const (
	WSS_URL_TARGET = "wss://socket.pasino.io/dice/"
	PORT           = 8080
)

// Handler is the exported function that Vercel will use as the entry point.
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/wss" {
		websocket.Handler(WebSocketHandler).ServeHTTP(w, r)
		return
	}
	http.NotFound(w, r)
}

// WebSocketHandler handles the WebSocket connection.
func WebSocketHandler(c *websocket.Conn) {
	ws := websocket.Message
	wssTarget, err := websocket.Dial(WSS_URL_TARGET, "", "http://localhost/")
	if err != nil {
		log.Println(err)
		if err := ws.Send(c, err.Error()); err != nil {
			panic(err)
		}
		return
	}

	defer wssTarget.Close()
	defer c.Close()
	go CopyMessages(&ws, c, wssTarget)
	go CopyMessages(&ws, wssTarget, c)

	select {}
}

func CopyMessages(ws *websocket.Codec, src, dst *websocket.Conn) {
	defer dst.Close()

	for {
		srcMsg := ""
		if err := ws.Receive(src, &srcMsg); err == io.EOF {
			log.Println("client disconnected")
			_ = ws.Send(src, err.Error())
			return
		} else if err != nil {
			log.Println(err)
			_ = ws.Send(src, err.Error())
			return
		}

		if err := ws.Send(dst, srcMsg); err != nil {
			log.Println("failed to send to websocket external")
			_ = ws.Send(src, err.Error())
			return
		}
	}
}

func main() {
	http.HandleFunc("/", Handler)
	log.Printf("starting server at port %v\n", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", PORT), nil))
}
