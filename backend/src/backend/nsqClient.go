package main

import (
	"fmt"
	"time"

	"github.com/bitly/go-nsq"
)

type NsqClient struct{}

var client = NsqClient{}

func (k *NsqClient) run(h *hub, nsqLookupds []string) {
	var reader *nsq.Consumer
	var err error
	inChan := make(chan *nsq.Message)
	conf := nsq.NewConfig()
	conf.Set("maxInFlight", 1000)
	reader, err = nsq.NewConsumer("all", "testqueue#ephemeral", conf)
	if err != nil {
		fmt.Printf("Error: %+v \n", err)
		select {
		case <-time.After(10 * time.Second):
			go k.run(h, nsqLookupds)
		}
		return
	}

	reader.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		inChan <- m
		return nil
	}))
	err = reader.ConnectToNSQLookupds(nsqLookupds)

	if err != nil {
		fmt.Printf("Error: %+v \n", err)
		select {
		case <-time.After(10 * time.Second):
			go k.run(h, nsqLookupds)
		}
		return
	}

	go loop(h, inChan)
}

func loop(h *hub, inChan chan *nsq.Message) {
	for msg := range inChan {
		fmt.Printf("got msg: %+v\n", string(msg.Body))
		h.broadcast <- msg.Body
		msg.Finish()
	}
}
