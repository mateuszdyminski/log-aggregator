package main

import (
	"fmt"
	"os"
	"time"

	"github.com/bitly/go-nsq"
	"github.com/mateuszdyminski/glog"
)

func SetupLogNsq(nsqHost string) {

	cfg := nsq.NewConfig()

	// make the producer
	producer, err := nsq.NewProducer(nsqHost+":4150", cfg)
	if err != nil {
		fmt.Printf("Error: %+v \n", err)
		select {
		case <-time.After(10 * time.Second):
			go SetupLogNsq(nsqHost)
		}
		return
	} else {
		fmt.Println("> connected")
	}

	host, err := os.Hostname()
	if err != nil {
		host = "localhost"
	}

	logsChannel := make(chan []byte, 1000)
	glog.SetChannel(logsChannel)

	go func(nsqHost, host string) {
		for {
			select {
			case log := <-logsChannel:
				if err = producer.Publish(host, log); err != nil {
					fmt.Printf("Error: %+v \n", err)
				}
				if err = producer.Publish("all", log); err != nil {
					fmt.Printf("Error: %+v \n", err)
				}

				fmt.Println("Logs sent\n")
			}
		}
	}(nsqHost, host)
}

func Info(args ...interface{}) {
	glog.InfoDepth(1, args...)
}

func Infof(format string, args ...interface{}) {
	glog.InfoDepth(1, fmt.Sprintf(format, args...))
}

func Warning(args ...interface{}) {
	glog.WarningDepth(1, args...)
}

func Warningf(format string, args ...interface{}) {
	glog.WarningDepth(1, fmt.Sprintf(format, args...))
}

func Error(args ...interface{}) {
	glog.ErrorDepth(1, args...)
}

func Errorf(format string, args ...interface{}) {
	glog.ErrorDepth(1, fmt.Sprintf(format, args...))
}

func Fatal(format string, args ...interface{}) {
	glog.FatalDepth(1, args...)
}

func Fatalf(format string, args ...interface{}) {
	glog.FatalDepth(1, fmt.Sprintf(format, args...))
}

func V(level glog.Level) glog.Verbose {
	return glog.V(level)
}
