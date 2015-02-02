// TO DO
// Integrate LiveReload to gin architecture.

package gin

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// A little hackish but it works...
	CheckOrigin: func(req *http.Request) bool {
		return true
	},
}

type Handshake struct {
	Command    string   `json:"command"`
	Protocols  []string `json:"protocols"`
	ServerName string   `json:"serverName"`
}

func serveScript(rw http.ResponseWriter, req *http.Request) {
	http.ServeFile(rw, req, "./static/livereload.js")
}

func writeHandshake(ws *websocket.Conn, hs Handshake) {
	var sv Handshake

	sv.Command = "hello"
	sv.Protocols = hs.Protocols
	sv.ServerName = "Go-LiveReload"

	p, err := json.Marshal(&sv)

	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Emmiting server handshake...")
	ws.WriteMessage(1, p)
}

func readWebSocket(ws *websocket.Conn) {
	defer ws.Close()

	for {
		_, p, err := ws.ReadMessage()

		if err != nil {
			log.Println(err)
			return
		}

		var hs Handshake

		err = json.Unmarshal(p, &hs)

		if err != nil {
			log.Println(err)
			return
		}

		// Handshake
		if hs.Command == "hello" {
			log.Println("Handshake received...")
			writeHandshake(ws, hs)
		}
	}
}

func serveWebSockets(rw http.ResponseWriter, req *http.Request) {
	ws, err := upgrader.Upgrade(rw, req, nil)

	if err != nil {
		log.Println(err)
		return
	}

	readWebSocket(ws)
}

func main() {
	log.Println("Initializing LiveReload server...")

	// Serving static livereload client script
	http.HandleFunc("/livereload.js", serveScript)

	// Serving sockets
	http.HandleFunc("/livereload", serveWebSockets)

	// Listen up!
	http.ListenAndServe(":35729", nil)
}
