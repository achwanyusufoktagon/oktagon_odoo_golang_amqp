package main

import (
	"github.com/streadway/amqp"
	"oktagon_odoo_golang_amqp/readModel"
	"log"
	"fmt"
)

func main() {
	// Declaring AMQP cloud and queue exchange name
	var amqpURI = "amqp://kjaexlqs:planetMars5050@feisty-peacock.rmq.cloudamqp.com:5672/kjaexlqs"
	var exchangeComsumer = "oktagon_odoo_amqp"
	var exchangeProvider = "oktagon_odoo_providing_amqp"
	log.Printf("dialing %q", amqpURI)
	// Dialing connection
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		fmt.Errorf("Dial: %s", err)
	}
	go func() {
		fmt.Printf("closing: %s", <-conn.NotifyClose(make(chan *amqp.Error)))
	}()
	// Connecting to channel
	log.Printf("got Connection, getting Channel")
	ch, err := conn.Channel()
	if err != nil {
		fmt.Errorf("Channel: %s", err)
	}
	// Declaring Exchange Consumer
	log.Printf("got Channel, declaring Exchange Consumer (%q)", exchangeComsumer)
	q, err := ch.QueueDeclare(
		exchangeComsumer,
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		fmt.Errorf("Error while declaring Exchange Consumer: %s", err)
	}
	// Declaring Exchange Provider
	log.Printf("got Channel, declaring Exchange Publisher (%q)", exchangeProvider)
	r, err := ch.QueueDeclare(
		exchangeProvider,
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		fmt.Errorf("Error while declaring Exchange Consumer: %s", err)
	}
	// Consuming Order
	log.Printf("starting Consume")
	msgs, _ := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	log.Printf("Processing Consumed Data")
	log.Println(readModel.GetAllProduct())
	for m := range msgs {
		n := string(m.Body)
		if n == "Get Product" {
			// Publishing Protobuf Data (All Product)
			log.Printf("Publishing Product Data")
			msg := amqp.Publishing{
				ContentType:   "text/plain",
                CorrelationId: m.CorrelationId,
				Body: readModel.GetAllProduct(),
			}
			ch.Publish(
				"",
				r.Name,
				false,
				false,
				msg)
		}
	}
}