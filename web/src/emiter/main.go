package main

import (
	"flag"
	"fmt"
	"html"
	"net/http"
	"time"
)

var (
	nsqHost string
	port    int
)

func init() {
	// Flags
	flag.Usage = func() {
		flag.PrintDefaults()
	}

	flag.StringVar(&nsqHost, "nsqHost", "nsq-worker", "nsq host")
	flag.IntVar(&port, "p", 8001, "web port")
}

func main() {
	flag.Parse()
	fmt.Printf("Nsq host: %s\n", nsqHost)
	fmt.Printf("Web port: %d\n", port)

	SetupLogNsq(nsqHost)

	http.HandleFunc("/", dummyHandler)
	fmt.Printf("Web port 1")
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	fmt.Printf("Web port 2")
}

func dummyHandler(rw http.ResponseWriter, req *http.Request) {
	Infof("%s - - [%s] \"%s %s %s\" \"%s\" \"%s\"\n", req.RemoteAddr,
		time.Now().Format("02/Jan/2006:15:04:05 -0700"), req.Method, req.RequestURI, req.Proto, req.Referer(), req.UserAgent())
	fmt.Fprintf(rw, "Hello, %q", html.EscapeString(req.URL.Path))
}
