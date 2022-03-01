package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func main() {
	stuck := make(chan os.Signal)

	ServeHTTP()

	<-stuck
}

// ServeHTTP ServeHTTP
func ServeHTTP() {
	go func() {
		g := gin.New()
		g.Use(gin.Recovery())
		err := g.SetTrustedProxies(nil)
		if err != nil {
			panic(err)
		}

		public := g.Group("/socket")
		public.GET("", SocketHandler)

		if err := g.Run(":8080"); err != nil {
			panic(err)
		}
	}()
}

func SocketHandler(c *gin.Context) {
	upGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}

	defer func() {
		closeSocketErr := ws.Close()
		if closeSocketErr != nil {
			panic(err)
		}
	}()

	for {
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Message Type: %d, Message: %s\n", msgType, string(msg))
		err = ws.WriteJSON(struct {
			Reply string `json:"reply"`
		}{
			Reply: "Echo...",
		})
		if err != nil {
			panic(err)
		}
	}
}
