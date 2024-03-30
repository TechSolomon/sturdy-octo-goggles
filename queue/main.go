package main

import (
	"fmt"
	"math"
	"sync" // For Lock & Unlock
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Tree struct {
	Left  *Tree
	Value int
	Right *Tree
}

type SafeCounter struct {
	mu sync.Mutex
	v  map[string]int
}

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

func summation(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum
}

func subscribe(client mqtt.Client) {
	topic := "example/message"
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf(">> 🔔 Subscribed to topic: %s", topic)
}

func (c *SafeCounter) Inc(key string) {
	c.mu.Lock()
	c.v[key]++
	c.mu.Unlock()
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

func (c *SafeCounter) Value(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.v[key]
}

func publish(client mqtt.Client) {
	example := 42.0
	solution := intermediary(example)
	context := fmt.Sprintf("Hello, world. %e", solution)
	token := client.Publish("example/message", 0, false, context)
	token.Wait()
	time.Sleep(time.Second)
}

func example() {
	c := SafeCounter{v: make(map[string]int)}
	for i := 0; i < SAMPLE; i++ {
		go c.Inc("🔑")
	}

	time.Sleep(time.Second)
	fmt.Println(c.Value("🔑"))
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

	fibo := []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144} // ! Range & Close

	c := make(chan int) // ? Buffered Channels
	go summation(fibo[:len(fibo)/2], c)
	go summation(fibo[len(fibo)/2:], c)
	x, y := <-c, <-c

	synchronization("second")

	publish(client)

	fmt.Println(x, y, x+y)

	example()
}
