package main

import (
	"github.com/googollee/go-socket.io"
	"log"
	"net/http"
	"time"
)

func main() {
	accountsId := []string{}
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")

		so.Emit("connection", "success")

		so.On("accountId", func(msg string) {
			so.Join(msg)
			accountsId = append(accountsId, msg)
		})

		so.On("disconnection", func() {
			log.Println("on disconnect")
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	time.AfterFunc(5*time.Second, func() {
		log.Println("timer")
		sendMessage(server, accountsId[0], "this is goooo!")
	})

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:4444...")
	log.Fatal(http.ListenAndServe(":4444", nil))
}

func sendMessage(server *socketio.Server, accountId string, message string) {
	server.BroadcastTo(accountId, "notify", message)
}
