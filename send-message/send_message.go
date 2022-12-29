package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Israel-Ferreira/alura-sqs/common"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

const QUEUE_NAME string = "alura-teste"

func SendMessage(sess *session.Session, queueUrl string) error {
	svc := sqs.New(sess)

	transfer := common.NovaTransferencia(
		common.Conta{Agencia: 00001, Conta: "00021122-3"},
		common.Conta{Agencia: 00004, Conta: "00002211-1"},
		500.00,
		"BRL",
	)

	jsonStr, err := json.Marshal(&transfer)

	if err != nil {
		return err
	}

	messageBody := string(jsonStr)

	_, err = svc.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    &queueUrl,
		MessageBody: &messageBody,
	})

	if err != nil {
		return err
	}

	return nil
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

	queueName := QUEUE_NAME

	result, err := common.GetQueueUrl(sess, &queueName)

	if err != nil {
		log.Fatalf("Erro ao obter a url da Fila: %v \n", err.Error())
	}

	fmt.Println(*result.QueueUrl)

	queueUrl := *result.QueueUrl

	if err = SendMessage(sess, queueUrl); err != nil {
		log.Fatalf("Erro ao enviar a mensagem: %v \n", err.Error())
	}

}
