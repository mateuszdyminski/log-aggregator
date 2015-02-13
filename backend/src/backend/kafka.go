package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

type KafkaClient struct{}

type Msg struct {
	Key     string
	Message string
}

var client = KafkaClient{}

func (k *KafkaClient) run(h *hub, kafkaTopic, kafkaHost string) {
	client, err := sarama.NewClient("my_client", []string{kafkaHost + ":9092"}, nil)
	fmt.Printf("Client: %+v \n", client)
	if err != nil {
		fmt.Printf("Error: %+v \n", err)
		select {
		case <-time.After(10 * time.Second):
			go k.run(h, kafkaTopic, kafkaHost)
		}
		return
	} else {
		fmt.Println("> connected")
	}
	defer client.Close()

	consumer, err := sarama.NewConsumer(client, kafkaTopic, 0, "group", sarama.NewConsumerConfig())
	if err != nil {
		fmt.Printf("Topic: %s, Error: %+v \n", kafkaTopic, err)
		select {
		case <-time.After(10 * time.Second):
			go k.run(h, kafkaTopic, kafkaHost)
		}
		return
	} else {
		fmt.Println("> consumer ready")
	}
	defer consumer.Close()

consumerLoop:
	for {
		select {
		case event := <-consumer.Events():
			if event.Err != nil {
				fmt.Printf("%+v \n", event.Err)
			} else {
				msg := Msg{string(event.Key), string(event.Value)}

				fmt.Printf("got msg: %+v\n", msg)

				data, err := json.Marshal(msg)
				if err != nil {
					fmt.Printf("%+v \n", err)
					break consumerLoop
				}

				h.broadcast <- data
			}
		case <-time.After(60 * time.Minute):
			fmt.Println("> timed out")
			break consumerLoop
		}
	}
}
