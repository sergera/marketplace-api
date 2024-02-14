package websocket

import (
	"log"
	"net/http"
	"strings"

	"github.com/sergera/marketplace-api/internal/conf"
	"nhooyr.io/websocket"
)

var wsOptions = websocket.AcceptOptions{
	InsecureSkipVerify: true,
	OriginPatterns:     strings.Split(conf.GetConf().CORSAllowedURLs, ","),
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	ws, err := websocket.Accept(w, r, &wsOptions)
	if err != nil {
		log.Println(err)
		return ws, err
	}

	return ws, nil
}
