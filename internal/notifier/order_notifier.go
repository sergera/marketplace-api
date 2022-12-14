package notifier

import (
	"context"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/sergera/marketplace-api/internal/conf"
	"github.com/sergera/marketplace-api/internal/domain"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

var once sync.Once
var instance *OrderNotifier

var mu sync.Mutex

type OrderNotifier struct {
	queue     []domain.OrderModel
	wsOptions websocket.AcceptOptions
}

func GetOrderNotifier() *OrderNotifier {
	once.Do(func() {
		conf := conf.GetConf()
		var n *OrderNotifier = &OrderNotifier{
			[]domain.OrderModel{},
			websocket.AcceptOptions{
				InsecureSkipVerify: true,
				OriginPatterns:     strings.Split(conf.CORSAllowedURLs, ","),
			},
		}
		instance = n
	})
	return instance
}

func (n *OrderNotifier) PushOrders(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &n.wsOptions)
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "oops, something went wrong")

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	ctx = c.CloseRead(ctx)

	t := time.NewTicker(time.Second)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			c.Close(websocket.StatusNormalClosure, "")
			return
		case <-t.C:
			if len(n.queue) > 0 {
				mu.Lock()
				if err := wsjson.Write(ctx, c, n.queue); err != nil {
					log.Println("error sending websocket message: ", err.Error())
				}
				n.queue = []domain.OrderModel{}
				mu.Unlock()
			}
		}
	}
}

func (n *OrderNotifier) AppendOrder(m domain.OrderModel) {
	mu.Lock()
	n.queue = append(n.queue, m)
	mu.Unlock()
}
