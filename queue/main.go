package main

import (
	"fmt"
	"math"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	SAMPLE = 42
)

var messagePublishHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf(">> Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf(">> ❌ Connection lost: %v", err)
}

var connHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println(">> ✅ Connection successful!")
}

func subscribe(client mqtt.Client) {
	topic := "example/message"
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf(">> 🔔 Subscribed to topic: %s", topic)
}

func intermediary(input float64) float64 {
	output := math.Sin(input)
	return output
}

func synchronization(s string) {
	for i := 0; i < SAMPLE; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
}

func publish(client mqtt.Client) {
	example := 42.0
	solution := intermediary(example)
	context := fmt.Sprintf("Hello, world. %e", solution)
	token := client.Publish("example/message", 0, false, context)
	token.Wait()
	time.Sleep(time.Second)
}

func main() {
	var broker, port = "localhost", 1883
	location := fmt.Sprintf("tcp://%s:%d", broker, port)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(location)

	opts.SetClientID("go_mqtt_sandbox")
	opts.SetDefaultPublishHandler(messagePublishHandler)

	opts.OnConnectionLost = connLostHandler
	opts.OnConnect = connHandler

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	subscribe(client)

	go synchronization("first")
	synchronization("second")

	publish(client)
}
