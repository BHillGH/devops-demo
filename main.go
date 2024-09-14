package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveSocket() http.Handler {
	hub := newHub()
	go hub.run()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

}

func main() {
	flag.Parse()

	dir := http.Dir("./frontend")
	fs := http.FileServer(dir)

	mux := http.NewServeMux()

	mux.Handle("/", fs)
	mux.Handle("/ws", serveSocket())

	err := http.ListenAndServe(*addr, mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
