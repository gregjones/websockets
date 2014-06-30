package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
)

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws}
	h.register <- c
	defer func() { h.unregister <- c }()
	go c.writer()
	c.reader()
}

var (
	addr      = flag.String("add", ":5432", "http addr")
	homeTempl *template.Template
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	homeTempl.Execute(w, r.Host)
}

func main() {
	flag.Parse()
	homeTempl = template.Must(template.ParseFiles("index.html"))
	go h.run()
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ws", wsHandler)
	log.Printf("Listening on http://localhost%s", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("Failed to serve:", err)
	}
}
