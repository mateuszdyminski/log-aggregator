package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bitly/go-nsq"
)

type NsqClient struct{}

type Msg struct {
	Key     string
	Message string
}

var client = NsqClient{}

func (k *NsqClient) run(h *hub, nsqTopic, nsqHost string) {
	var reader *nsq.Consumer
	var err error
	inChan := make(chan *nsq.Message)
	lookup := nsqHost + ":4161"
	conf := nsq.NewConfig()
	conf.Set("maxInFlight", 1000)
	reader, err = nsq.NewConsumer(nsqTopic, "testqueue#ephemeral", conf)
	if err != nil {
		fmt.Printf("Error: %+v \n", err)
		select {
		case <-time.After(10 * time.Second):
			go k.run(h, nsqTopic, nsqHost)
		}
		return
	}

	reader.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		inChan <- m
		return nil
	}))
	err = reader.ConnectToNSQLookupd(lookup)

	if err != nil {
		fmt.Printf("Error: %+v \n", err)
		select {
		case <-time.After(10 * time.Second):
			go k.run(h, nsqTopic, nsqHost)
		}
		return
	}

	go loop(h, inChan)
}

func loop(h *hub, inChan chan *nsq.Message) {
	for msg := range inChan {
		message := Msg{"all", string(msg.Body)}

		fmt.Printf("got msg: %+v\n", message)

		data, err := json.Marshal(message)
		if err != nil {
			fmt.Printf("%+v \n", err)
		} else {
			h.broadcast <- data
			msg.Finish()
		}
	}
}
