package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/getsentry/sentry-go"
	_ "github.com/joho/godotenv/autoload"
	"go-kafka-simple/variables"
	"os"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(variables.NumCore)
	if err := sentry.Init(sentry.ClientOptions{
		Dsn: variables.SentryLink,
	}); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to init sentry: %s\n", err)
		panic(err)
	}
}

func main() {
	fmt.Println("Consumer is starting...")

	conf := kafka.ConfigMap{
		"bootstrap.servers":  variables.KafkaBootstrapServers,
		"group.id":           variables.KafkaGroupId,
		"auto.offset.reset":  "earliest",
	}

	consumer, err := kafka.NewConsumer(&conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		panic(err)
	}

	consumer.SubscribeTopics([]string{variables.KafkaTopic}, nil)

	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			sentry.CaptureException(err)
		}
	}

	consumer.Close()

	fmt.Println("Consumer done")
}
