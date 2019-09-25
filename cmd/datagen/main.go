package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/ahammar-dkj/datagen/internal/data"
	"github.com/ahammar-dkj/datagen/internal/sqs"
)

var (
	queueUrl string
	region   string
	endpoint string
	interval int64
)

func init() {
	flag.StringVar(&queueUrl, "queue-url", "", "AWS SQS queue URL")
	flag.StringVar(&region, "region", "us-east-1", "AWS Region")
	flag.StringVar(&endpoint, "endpoint", "", "AWS SQS service endpoint")
	flag.Int64Var(&interval, "interval", 1000, "Interval (ms) between to messages")
	flag.Parse()
}

func main() {
	if queueUrl == "" {
		flag.Usage()
		os.Exit(1)
	}

	sqs := sqs.NewClient(queueUrl, endpoint, region)
	log.Printf("Initializing new SQS client, queue URL: %s\n", queueUrl)

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		done <- true
	}()

	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	fmt.Printf("Sending messages, interval between messages: %d ms\n", interval)
	msgCount := 0

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				_, err := sqs.SendMessage(generateEvent())
				if err != nil {
					log.Printf(fmt.Sprintf("Failed to send message: %v\n", err))
				}
				msgCount++
				if msgCount%100 == 0 {
					log.Printf("%d messages sent\n", msgCount)
				}
			}
		}
	}()

	<-done
	ticker.Stop()

	log.Println("Done")
}

func generateEvent() interface{} {
	l, err := data.NewLogin()
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to create new login: %v\n", err))
	}

	extensions := map[string]interface{}{"subject": fmt.Sprintf("user/%s", l.UserName)}
	data := struct {
		Id        string  `json:"id"`
		Email     string  `json:"email"`
		Latitude  float32 `json:"latitude"`
		Longitude float32 `json:"longitude"`
	}{Id: l.Id, Email: l.Email, Latitude: l.Position.Latitude, Longitude: l.Position.Longitude}

	return struct {
		CloudEventsVersion string                 `json:"cloudEventsVersion,omitempty"`
		EventID            string                 `json:"eventID"`
		EventTime          time.Time              `json:"eventTime,omitempty"`
		EventType          string                 `json:"eventType"`
		ContentType        string                 `json:"contentType,omitempty"`
		Source             string                 `json:"source"`
		Extensions         map[string]interface{} `json:"extensions,omitempty"`
		Data               interface{}            `json:"data,omitempty"`
	}{
		CloudEventsVersion: "0.1",
		EventID:            uuid.New().String(),
		EventTime:          time.Now().UTC(),
		EventType:          "com.fake.datagen.login",
		ContentType:        "application/json",
		Source:             "com.fake/datagen",
		Extensions:         extensions,
		Data:               data,
	}
}
