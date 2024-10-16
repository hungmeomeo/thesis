package main

import (
	"demo/ws"
	"fmt"
	"log"
	"net/http"

	_ "net/http/pprof"
)

func main() {

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	hub := ws.NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.HandleWebSocket(hub, w, r)
	})

	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
