package websocket

type Pool struct {
	Register      chan *Connection
	Unregister    chan *Connection
	Connections   map[*Connection]bool
	BroadcastJSON chan interface{}
}

func NewPool() *Pool {
	p := &Pool{
		Register:      make(chan *Connection),
		Unregister:    make(chan *Connection),
		Connections:   make(map[*Connection]bool),
		BroadcastJSON: make(chan interface{}),
	}
	go p.Start()
	return p
}

func (pool *Pool) Start() {
	for {
		select {
		case conn := <-pool.Register:
			pool.Connections[conn] = true
		case conn := <-pool.Unregister:
			delete(pool.Connections, conn)
		case msg := <-pool.BroadcastJSON:
			for conn := range pool.Connections {
				go conn.WriteJSON(msg)
			}
		}
	}
}
