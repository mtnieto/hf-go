package main

import (
	"github.com/coren-hfservice/api/util/log"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/events/client"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/streadway/amqp"
)

const (
	chainCodeID = "coren"
	eventName   = "createAsset2"
)

func main() {
	// ConfigBackend contains config backend for integration tests

	// RABBIT MQ CONNECTION
	conn, err := amqp.Dial("amqp://user:bitnami@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	sdk, err := fabsdk.New(config.FromFile("../config/config-hf-dev.yaml"))
	if err != nil {
		log.Errorf("[prepareSDK] Failed to create new SDK: ", err)
	}

	ctxProvider := sdk.ChannelContext("channel", fabsdk.WithUser("admin"), fabsdk.WithOrg("org1"))

	chContext, err := ctxProvider()
	if err != nil {
		log.Errorf("[prepareSDK] Failed to get channel context: ", err)
	}

	eventService, err := chContext.ChannelService().EventService(client.WithBlockEvents())
	_, cceventch, err := eventService.RegisterChaincodeEvent(chainCodeID, "createAsset2")
	if err != nil {
		log.Fatalf("Error registering for filtered block events: %s", err)
	}
	// _, beventch, err := eventService.RegisterBlockEvent()
	// if err != nil {
	// 	log.Errorf("Error registering for block events: %s", err)
	// }
	event, ok := <-cceventch
	log.Infof("LLega aqui")
	log.Infof("%v", ok)
	log.Infof("%v", cceventch)
	log.Infof("%v", string(event.Payload[:]))

	body := string(event.Payload[:])
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	// log.Infof("%v", event.Payload)
}
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
