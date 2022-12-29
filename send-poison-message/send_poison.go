package main

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/Israel-Ferreira/alura-sqs/common"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func SendPoison(sess *session.Session, queueUrl string) error {
	svc := sqs.New(sess)

	poisonContent := `
		<?xml version="1.0" encoding="UTF-8" ?>
		<root>
		<conta_origem>
			<agencia>%d</agencia>
			<numero_conta>%s</numero_conta>
		</conta_origem>
		<conta_destino>
			<agencia>%d</agencia>
			<numero_conta>%s</numero_conta>
		</conta_destino>
		<valor>%.2f</valor>
		<moeda>BRL</moeda>
		</root>
	`

	poisonFormatted := fmt.Sprintf(poisonContent, 05325, "099999-2", 0233, "000211-2", 500.00)

	xmlContent, err := xml.Marshal(poisonFormatted)

	if err != nil {
		return err
	}

	content := string(xmlContent)

	_, err = svc.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    &queueUrl,
		MessageBody: &content,
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

	queueName := "alura-teste"

	result, err := common.GetQueueUrl(sess, &queueName)

	if err != nil {
		log.Fatalf("Erro ao obter a url da Fila: %v \n", err.Error())
	}

	queueUrl := *result.QueueUrl

	if err := SendPoison(sess, queueUrl); err != nil {
		log.Fatalln("Erro ao mandar a mensagem com defeito")
	}

}
