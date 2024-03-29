package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/sergera/marketplace-api/internal/conf"
	"github.com/sergera/marketplace-api/internal/domain"
	"github.com/sergera/marketplace-api/internal/evt"
	"github.com/sergera/marketplace-api/internal/notifier"
	"github.com/sergera/marketplace-api/internal/repositories"
)

type OrderAPI struct {
	orderRepo    *repositories.OrderRepository
	eventHandler *evt.EventHandler
	notifier     *notifier.OrderNotifier
}

func NewOrderAPI() *OrderAPI {
	conf := conf.GetConf()
	conn := repositories.NewDBConnection(conf.DBHost, conf.DBPort, conf.DBName, conf.DBUser, conf.DBPassword, false)
	conn.Open()
	orderRepo := repositories.NewOrderRepository(conn)
	evtHandler := evt.NewEventHandler()
	orderNotifier := notifier.GetOrderNotifier()
	return &OrderAPI{orderRepo, evtHandler, orderNotifier}
}

func (o *OrderAPI) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var m domain.OrderModel

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := m.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m.Date = time.Now()

	if err := o.orderRepo.CreateOrder(&m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderInBytes, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	o.eventHandler.Produce(evt.Topics[domain.Unconfirmed], "", orderInBytes)
	o.notifier.Publish(m)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(orderInBytes)
}

func (o *OrderAPI) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	var m domain.OrderModel

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := m.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := o.orderRepo.UpdateOrder(m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	o.notifier.Publish(m)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}

func (o *OrderAPI) GetOrderRange(w http.ResponseWriter, r *http.Request) {
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	oldestFirst := r.URL.Query().Get("oldest-first")

	if oldestFirst == "" {
		oldestFirst = "false"
	}

	oldestFirstBool, err := strconv.ParseBool(oldestFirst)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m := domain.OrderRangeModel{
		Start:       start,
		End:         end,
		OldestFirst: oldestFirstBool,
	}

	if err := m.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orders, err := o.orderRepo.GetOrderRange(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ordersInBytes, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(ordersInBytes)
}

func CorsHandler(h http.HandlerFunc) http.HandlerFunc {
	conf := conf.GetConf()
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", conf.CORSAllowedURLs)
		if r.Method == "OPTIONS" {
			//handle preflight in here
			w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Accept")
			w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		} else {
			h.ServeHTTP(w, r)
		}
	}
}
