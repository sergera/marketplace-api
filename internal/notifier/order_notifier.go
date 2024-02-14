package notifier

import (
	"log"
	"net/http"
	"sync"

	"github.com/sergera/marketplace-api/pkg/websocket"
)

var once sync.Once
var instance *OrderNotifier

type OrderNotifier struct {
	pool *websocket.Pool
}

func GetOrderNotifier() *OrderNotifier {
	once.Do(func() {
		pool := websocket.NewPool()
		instance = &OrderNotifier{pool}
	})
	return instance
}

func (n *OrderNotifier) Subscribe(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r)
	if err != nil {
		log.Println(err)
	}

	n.pool.Register <- websocket.NewConnection(ws, n.pool)
}

func (n *OrderNotifier) Publish(msg interface{}) {
	n.pool.BroadcastJSON <- msg
}
