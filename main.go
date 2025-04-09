package main

import (
	"crypto/sha512"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
	"github.com/xdg-go/scram"
	"log"
	"os"
	"os/signal"
	"strings"
)

var (
	SHA512 scram.HashGeneratorFcn = sha512.New
)

type XDGSCRAMClient struct {
	*scram.Client
	*scram.ClientConversation
	scram.HashGeneratorFcn
}

func (x *XDGSCRAMClient) Begin(userName, password, authzID string) (err error) {
	x.Client, err = x.HashGeneratorFcn.NewClient(userName, password, authzID)
	if err != nil {
		return err
	}
	x.ClientConversation = x.Client.NewConversation()
	return nil
}

func (x *XDGSCRAMClient) Step(challenge string) (response string, err error) {
	response, err = x.ClientConversation.Step(challenge)
	return
}

func (x *XDGSCRAMClient) Done() bool {
	return x.ClientConversation.Done()
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	splitBrokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Consumer.Return.Errors = true
	conf.Metadata.Full = true
	conf.Net.SASL.Enable = true
	conf.Net.SASL.User = os.Getenv("KAFKA_USERNAME")
	conf.Net.SASL.Password = os.Getenv("KAFKA_PASSWORD")
	conf.Net.SASL.Handshake = true
	conf.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA512} }
	conf.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512

	master, err := sarama.NewConsumer(splitBrokers, conf)
	if err != nil {
		log.Fatal("Error creating the consumer:", err)
	}

	defer func() {
		if errMaster := master.Close(); err != nil {
			panic(errMaster)
		}
	}()

	consumer, err := master.ConsumePartition(os.Getenv("KAFKA_TOPIC"), 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	msgCount := 0
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case errConsumer := <-consumer.Errors():
				fmt.Println(errConsumer)
			case msg := <-consumer.Messages():
				msgCount++
				fmt.Println("Received messages", string(msg.Key), string(msg.Value))
			case <-signals:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	fmt.Println("Processed", msgCount, "messages")
}
