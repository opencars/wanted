package main

import (
	"log"
	"testing"

	nats "github.com/nats-io/nats.go"
)

func BenchmarkNATS(b *testing.B) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		nc.Publish("gov-data", []byte("{ \"test\": \"success\"}"))
	}

	b.StopTimer()
}
