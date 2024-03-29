package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sergera/marketplace-api/internal/api"
	"github.com/sergera/marketplace-api/internal/conf"
	"github.com/sergera/marketplace-api/internal/notifier"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	conf := conf.GetConf()

	mux := http.NewServeMux()

	orderAPI := api.NewOrderAPI()
	mux.HandleFunc("/create", api.CorsHandler(orderAPI.CreateOrder))
	mux.HandleFunc("/order-range", api.CorsHandler(orderAPI.GetOrderRange))
	mux.HandleFunc("/update-order", orderAPI.UpdateOrder)

	orderNotifier := notifier.GetOrderNotifier()
	mux.HandleFunc("/notify-orders", api.CorsHandler(orderNotifier.Subscribe))

	srv := &http.Server{
		Addr:    ":" + conf.Port,
		Handler: mux,
	}

	fmt.Printf("starting application on port %s", conf.Port)

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
