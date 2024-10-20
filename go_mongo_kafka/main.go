package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// MongoDB connection setup
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}
	defer client.Disconnect(context.TODO())

	// Kafka setup
	brokerAddress := "localhost:9092"

	// Start a producer for testing
	go func() {
		producer(brokerAddress, "prism-user-create")
		producer(brokerAddress, "prism-user-delete")
		producer(brokerAddress, "prism-user-update")
	}()

	// Start a consumer for each topic
	go consumer(brokerAddress, "prism-user-create", client)
	go consumer(brokerAddress, "prism-user-delete", client)
	go consumer(brokerAddress, "prism-user-update", client)

	// Allow the consumers to run for some time
	time.Sleep(20 * time.Second)
}

// Producer sends messages to Kafka
func producer(brokerAddress, topic string) {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{brokerAddress},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})

	defer writer.Close()

	for i := 0; i < 5; i++ {
		msg := kafka.Message{
			Key:   []byte(fmt.Sprintf("Key-%d", i)),
			Value: []byte(fmt.Sprintf("This is message number %d", i)),
		}

		err := writer.WriteMessages(context.Background(), msg)
		if err != nil {
			log.Fatalf("could not write message %v", err)
		}

		fmt.Printf("Produced message: %s to topic %s\n", msg.Value, topic)
		time.Sleep(1 * time.Second)
	}
}

// Consumer reads messages from Kafka and writes them to MongoDB
func consumer(brokerAddress, topic string, mongoClient *mongo.Client) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{brokerAddress},
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	defer r.Close()

	fmt.Printf("Consumer started for topic: %s...\n", topic)

	// prism-user-create
	// prism-user-delete
	// prism-user-update
	var collection *mongo.Collection
	if topic == "prism-user-update" {
		collection = mongoClient.Database("mydb").Collection("user-update")
	} else if topic == "prism-user-create" {
		collection = mongoClient.Database("mydb").Collection("user-new")
	} else if topic == "prism-user-delete" {
		collection = mongoClient.Database("mydb").Collection("user-delete")
	} else {
		log.Fatalf("Unknown topic: %s", topic)
	}

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Fatalf("could not read message %v", err)
		}

		// Write message to MongoDB
		doc := bson.D{
			{Key: "key", Value: string(m.Key)},
			{Key: "value", Value: string(m.Value)},
			{Key: "timestamp", Value: time.Now()},
		}

		_, err = collection.InsertOne(context.TODO(), doc)
		if err != nil {
			log.Fatalf("could not insert document: %v", err)
		}

		fmt.Printf("Consumed message: key=%s value=%s from topic %s and saved to MongoDB\n", string(m.Key), string(m.Value), topic)
	}
}
