package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8000", "http server address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	flag.Parse()

	wsServer := NewWebsocketServer()
	go wsServer.Run()
	go listenTcp()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(wsServer, w, r)
	})

	http.HandleFunc("/", serveHome)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
