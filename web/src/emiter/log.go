package main

import (
	"fmt"
	"time"

	"github.com/bitly/go-nsq"
	"github.com/mateuszdyminski/glog"
)

func SetupLogNsq(nsqdAddr, host string) {

	cfg := nsq.NewConfig()

	// make the producer
	producer, err := nsq.NewProducer(nsqdAddr, cfg)
	if err != nil {
		fmt.Printf("Error: %+v \n", err)
		select {
		case <-time.After(10 * time.Second):
			go SetupLogNsq(nsqdAddr, host)
		}
		return
	} else {
		fmt.Println("> connected")
	}

	logsChannel := make(chan []byte, 1000)
	glog.SetChannel(logsChannel)

	go func(host string) {
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
	}(host)
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
