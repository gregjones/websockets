package main

type hub struct {
	connections map[*connection]bool
	broadcast   chan []byte
	register    chan *connection
	unregister  chan *connection
}

var h = hub{
	broadcast:   make(chan []byte),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

func (h *hub) doBroadcast(m []byte) {
	for c := range h.connections {
		h.send(c, m)
	}
}

func (h *hub) send(c *connection, m []byte) {
	select {
	case c.send <- m:
	default:
		delete(h.connections, c)
		close(c.send)
	}
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
			h.doBroadcast([]byte("Joined"))

		case c := <-h.unregister:
			delete(h.connections, c)
			h.doBroadcast([]byte("Left"))
			close(c.send)

		case m := <-h.broadcast:
			h.doBroadcast(m)
		}
	}
}
