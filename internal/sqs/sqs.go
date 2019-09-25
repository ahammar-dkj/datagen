package sqs

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Client struct {
	sqs      *sqs.SQS
	queueUrl string
}

func NewClient(queueUrl, endpoint, region string) *Client {
	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint: aws.String(endpoint),
		Region:   aws.String(region),
	}))
	sqs := sqs.New(sess)

	return &Client{sqs: sqs, queueUrl: queueUrl}
}

func (c *Client) SendMessage(message interface{}) (*string, error) {
	bytes, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	output, err := c.sqs.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(string(bytes)),
		QueueUrl:    &c.queueUrl,
	})
	if err != nil {
		return nil, err
	}
	return output.MessageId, nil
}
