package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var (
	hostname string
	nsqHost  string
	nsqTopic string
	port     int
	waitTime int
)

func init() {
	// Flags
	flag.Usage = func() {
		flag.PrintDefaults()
	}

	flag.StringVar(&hostname, "h", "localhost", "hostname")
	flag.IntVar(&port, "p", 8080, "port")
	flag.StringVar(&nsqHost, "nsqHost", "nsq-master", "nsq host")
	flag.StringVar(&nsqTopic, "nsqTopic", "all", "nsq topic")
	flag.IntVar(&waitTime, "wait", 0, "wait time")
}

func main() {
	flag.Parse()

	fmt.Printf("NsqHost: %s, topic: %s\n", nsqHost, nsqTopic)

	// run websocket hub
	go h.run()

	go client.run(&h, nsqTopic, nsqHost)

	http.HandleFunc("/wsapi/ws", serveWs)
	fmt.Printf("Server started, host: %s, port: %d\n", hostname, port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}

// serverWs handles websocket requests from the peer.
func serveWs(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Registering client to WS")
	if req.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	ws, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		fmt.Printf("Error %+v\n", err)
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws}
	h.register <- c
	go c.writePump()
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
