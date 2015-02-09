package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var (
	staticDir string
	hostname  string
	port      int
)

func init() {
	// Flags
	flag.Usage = func() {
		flag.PrintDefaults()
	}

	flag.StringVar(&hostname, "h", "localhost", "hostname")
	flag.StringVar(&staticDir, "dir", "app", "app directory")
	flag.IntVar(&port, "p", 9001, "port")
}

func main() {

	flag.Parse()

	fmt.Printf("Server started, host: %s, port: %d, serving dir: %s\n", hostname, port, staticDir)

	r := mux.NewRouter()

	r.Handle("/{path:.*}", http.FileServer(http.Dir(staticDir)))
	http.Handle("/", &loggingHandler{r})

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}

type loggingHandler struct {
	http.Handler
}

func (h loggingHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	t := time.Now()
	h.Handler.ServeHTTP(w, req)

	elapsed := time.Since(t)
	fmt.Printf("%s [%s] \"%s %s %s\" \"%s\" \"%s\" \"Took: %s\"\n", req.RemoteAddr,
		t.Format("02/Jan/2006:15:04:05 -0700"), req.Method, req.RequestURI, req.Proto, req.Referer(), req.UserAgent(), elapsed)
}
