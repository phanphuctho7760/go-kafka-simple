package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go-kafka-simple/variables"
	"os"
)

func main() {
	fmt.Println("Producer is starting...")

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": variables.KafkaBootstrapServers})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create producer: %s\n", err)
		panic(err)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	for _, word := range []string{"This", "is", "go", "simple", "kafka"} {
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &variables.KafkaTopic, Partition: kafka.PartitionAny},
			Value:          []byte(word),
		}, nil)
	}

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)

	fmt.Println("Producer was done")
}
