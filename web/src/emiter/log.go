package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bitly/go-nsq"
	"github.com/mateuszdyminski/glog"
)

type LogMsg struct {
	Host    string
	Level   string
	Content string
}

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
				content := string(log)
				logData, logErr := json.Marshal(LogMsg{host, string(content[0]), content})
				if logErr != nil {
					fmt.Printf("%+v \n", logErr)
				} else {
					if err = producer.Publish(host, logData); err != nil {
						fmt.Printf("Error: %+v \n", err)
					}
					if err = producer.Publish("all", logData); err != nil {
						fmt.Printf("Error: %+v \n", err)
					}

					fmt.Println("Logs sent from host", host, nsqdAddr)
				}
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
