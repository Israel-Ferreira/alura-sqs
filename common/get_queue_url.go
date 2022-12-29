package common

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func GetQueueUrl(sess *session.Session, name *string) (*sqs.GetQueueUrlOutput, error) {
	client := sqs.New(sess)

	result, err := client.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: name,
	})

	if err != nil {
		return nil, err
	}

	return result, nil

}
