package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Teeworlds-Server-Moderation/common/events"
	"github.com/Teeworlds-Server-Moderation/common/mqtt"
)

// Connect to the broker and publish a message periodically

var (
	topic         = events.TypePlayerJoin
	serverAddress = "tcp://localhost:1883"
	clientID      = "pubsub"
)

func init() {
	if brokerAddress := os.Getenv("BROKER_ADDRESS"); brokerAddress != "" {
		serverAddress = brokerAddress
	}
	if id := os.Getenv("BROKER_CLIENT_ID"); id != "" {
		clientID = id
	}
	if t := os.Getenv("BROKER_TOPIC"); t != "" {
		topic = t
	}

	log.Println("Initialized with address: ", serverAddress, " clientID: ", clientID)

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

func main() {
	// Messages will be delivered asynchronously so we just need to wait for a signal to shutdown
	sig := make(chan os.Signal, 1)
	publisher, err := mqtt.NewPublisher(serverAddress, "pubsub-publisher", "default")
	if err != nil {
		log.Fatalln("Could not create Publisher:", err)
	}
	defer publisher.Close()

	subscriber, err := mqtt.NewSubscriber(serverAddress, "pubsub-subscriber", "default", "different", topic)
	if err != nil {
		log.Fatalln("Could not create Subscriber:", err)
	}
	defer subscriber.Close()

	go func() {
		cnt := 0
		for {
			select {
			case <-time.After(time.Second * 30):
				if cnt%2 == 0 {
					publisher.Publish(fmt.Sprintf("%d default", cnt))
				} else {
					publisher.PublishTo("different", fmt.Sprintf("%d different", cnt))
				}
				cnt++
			case <-sig:
				return
			}
		}
	}()

	go func() {
		log.Println("Started subscriber routine.")
		for msg := range subscriber.Next() {
			switch msg.Topic {
			case "different":
				log.Println("Received message(", msg.Topic, "): ", msg.Payload)
			default:
				log.Println("Received message(", msg.Topic, "): ", msg.Payload)
			}
		}
	}()
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	fmt.Println("Connection is up, press Ctrl-C to shutdown")
	<-sig
	fmt.Println("signal caught - exiting")
	fmt.Println("shutdown complete")
}
