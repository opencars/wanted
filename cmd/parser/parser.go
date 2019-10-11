package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	nats "github.com/nats-io/nats.go"
)

func producer(amount int) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	flag.Parse()

	start := time.Now()

	for i := 0; i < amount; i++ {
		nc.Publish("gov-data", []byte(""))
	}

	period := time.Since(start)

	fmt.Println(amount, period)

	nc.Close()
}

func consumer(amount int) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()

	res := 0

	ch := make(chan *nats.Msg, 64)
	sub, _ := nc.ChanSubscribe("gov-data", ch)

	for range ch {
		res++
		if res == amount {
			break
		}
	}

	sub.Unsubscribe()

	period := time.Since(start)

	fmt.Println(res, period)

	nc.Close()
}

func main() {
	flag.Parse()

	amount := flag.Int("amount", 10000000, "Amount of messages to be sent")

	go consumer(*amount)
	producer(*amount)
}
