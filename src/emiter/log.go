package main

import (
	"fmt"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/mateuszdyminski/glog"
)

type Config struct {
	Logs     chan []byte
	Brokers  []string
	Topic    string
	ClientId string
	Host     string
}

type KafkaProducer struct {
	c Config
}

type logger struct {
	p *KafkaProducer
}

var mu sync.Mutex
var l *logger
var p *KafkaProducer

func SetupLogToKafka() {
	mu.Lock()
	defer mu.Unlock()

	c := Config{
		Logs: make(chan []byte, 1000),
	}

	p := KafkaProducer{c}

	client, err := sarama.NewClient("client_id", []string{"localhost:9092"}, sarama.NewClientConfig())
	if err != nil {
		panic(err)
	} else {
		fmt.Println("> connected")
	}

	producer, err := sarama.NewProducer(client, nil)
	if err != nil {
		panic(err)
	}
	// defer producer.Close()

	glog.SetChannel(p.c.Logs)

	go func() {
		for {
			select {
			case log := <-p.c.Logs:
				producer.Input() <- &sarama.MessageToSend{Topic: "test", Key: sarama.StringEncoder("localhost"), Value: sarama.StringEncoder(log)}
			case err := <-producer.Errors():
				panic(err.Err)
			}
		}
	}()
}

func KafkaEnabled() bool {
	return glog.CustomChannel()
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
	mu.Lock()
	defer mu.Unlock()
}

func Errorf(format string, args ...interface{}) {
	glog.ErrorDepth(1, fmt.Sprintf(format, args...))
	mu.Lock()
	defer mu.Unlock()
}

func Fatal(format string, args ...interface{}) {
	glog.FatalDepth(1, args...)
	mu.Lock()
	defer mu.Unlock()
}

func Fatalf(format string, args ...interface{}) {
	glog.FatalDepth(1, fmt.Sprintf(format, args...))
	mu.Lock()
	defer mu.Unlock()
}

func V(level glog.Level) glog.Verbose {
	return glog.V(level)
}
