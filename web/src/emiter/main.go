package main

import (
	"flag"
	"fmt"
	"html"
	"net/http"
	"time"
)

var (
	kafka string
	topic string
	port  int
)

func init() {
	// Flags
	flag.Usage = func() {
		flag.PrintDefaults()
	}

	flag.StringVar(&kafka, "k", "kafka", "kafka host")
	flag.StringVar(&topic, "topic", "test", "kafka topic")
	flag.IntVar(&port, "p", 8001, "web port")
}

func main() {
	flag.Parse()
	fmt.Printf("Kafka host: %s\n", kafka)
	fmt.Printf("Kafka topic: %s\n", topic)
	fmt.Printf("Web port: %d\n", port)

	SetupLogToKafka(topic, kafka)

	http.HandleFunc("/", dummyHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

}

func dummyHandler(rw http.ResponseWriter, req *http.Request) {
	Infof("%s - - [%s] \"%s %s %s\" \"%s\" \"%s\"\n", req.RemoteAddr,
		time.Now().Format("02/Jan/2006:15:04:05 -0700"), req.Method, req.RequestURI, req.Proto, req.Referer(), req.UserAgent())
	fmt.Fprintf(rw, "Hello, %q", html.EscapeString(req.URL.Path))
}
