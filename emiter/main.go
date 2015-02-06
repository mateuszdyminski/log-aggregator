package main

import (
	"flag"
	"fmt"
	"html"
	"net/http"
	"time"

	"github.com/mateuszdyminski/glog"
)

func main() {
	flag.Parse()
	defer glog.Flush()
	glog.Info("Starting")
	http.HandleFunc("/", dummyHandler)
	glog.Fatal(http.ListenAndServe(":8001", nil))
}

func dummyHandler(rw http.ResponseWriter, req *http.Request) {
	glog.Infof("%s - - [%s] \"%s %s %s\" \"%s\" \"%s\"\n", req.RemoteAddr,
		time.Now().Format("02/Jan/2006:15:04:05 -0700"), req.Method, req.RequestURI, req.Proto, req.Referer(), req.UserAgent())
	fmt.Fprintf(rw, "Hello, %q", html.EscapeString(req.URL.Path))
}
