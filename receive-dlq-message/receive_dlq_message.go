package main

import (
	"fmt"
	"log"

	"github.com/Israel-Ferreira/alura-sqs/common"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func DeleteMessage(sess *session.Session, queueUrl string, msg *sqs.Message) error {
	svc := sqs.New(sess)

	_, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &queueUrl,
		ReceiptHandle: msg.ReceiptHandle,
	})

	if err != nil {
		return err
	}

	return nil
}

func ReceiveDlqMessages(sess *session.Session, queueUrl string, chMessage chan<- *sqs.Message) {
	svc := sqs.New(sess)

	for {
		msgs, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:          &queueUrl,
			WaitTimeSeconds:   aws.Int64(20),
			VisibilityTimeout: aws.Int64(30),
		})

		if err != nil {
			log.Println("Failed to get new Messages")
			continue
		}

		fmt.Println(len(msgs.Messages))

		for _, msg := range msgs.Messages {
			chMessage <- msg
		}

	}
}

func main() {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("sa-east-1"),
		},
		SharedConfigState: session.SharedConfigEnable,
	})

	if err != nil {
		log.Fatalf("Erro ao conectar com a AWS: %v \n", err.Error())
	}

	queueName := "alura-teste-dlq"

	result, err := common.GetQueueUrl(sess, &queueName)

	if err != nil {
		log.Fatalf("Erro ao obter a url da Fila: %v \n", err.Error())
	}

	queueUrl := *result.QueueUrl

	chMessages := make(chan *sqs.Message)

	go ReceiveDlqMessages(sess, queueUrl, chMessages)

	for msgMorta := range chMessages {
		fmt.Println(*msgMorta.Body)

		if err := DeleteMessage(sess, queueUrl, msgMorta); err != nil {
			log.Println(err)
			continue
		}

		log.Println("Mensagem processada com sucesso")

	}
}
