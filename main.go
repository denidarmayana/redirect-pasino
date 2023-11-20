package main

import (
        "fmt"
        "io"
        "log"
        "net/http"

        "golang.org/x/net/websocket"
)

const (
        WSS_URL_TARGET = "wss://socket.pasino.com/dice/"
        PORT           = 8080
)

func main() {
        server := http.DefaultServeMux

        server.Handle("/wss", websocket.Handler(func(c *websocket.Conn) {
                ws := websocket.Message
                wssTarget, err := websocket.Dial(WSS_URL_TARGET, "wss", "*")
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
        }))

        log.Printf("starting server at port %v\n", PORT)
        http.ListenAndServe(fmt.Sprintf(":%v", PORT), server)
}

func CopyMessages(ws *websocket.Codec, src, dst *websocket.Conn) {
        defer src.Close()
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