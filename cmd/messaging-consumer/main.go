package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/michaelchandrag/botfood-go/infrastructures/database"
	"github.com/michaelchandrag/botfood-go/internal/logger"
	"github.com/michaelchandrag/botfood-go/utils"

	dto "github.com/michaelchandrag/botfood-go/pkg/modules/messaging/dto"
	service "github.com/michaelchandrag/botfood-go/pkg/modules/messaging/services"
)

func GetQueueURL(sess *session.Session, queue string) (*sqs.GetQueueUrlOutput, error) {
	svc := sqs.New(sess)

	urlResult, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queue,
	})
	if err != nil {
		return nil, err
	}

	return urlResult, nil
}

func GetMessages(sess *session.Session, queueURL *string, timeout *int64) (*sqs.ReceiveMessageOutput, error) {
	svc := sqs.New(sess)

	msgResult, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            queueURL,
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   timeout,
	})
	if err != nil {
		return nil, err
	}

	return msgResult, nil
}

func DeleteMessage(sess *session.Session, queueURL *string, messageHandle *string) error {
	svc := sqs.New(sess)

	_, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      queueURL,
		ReceiptHandle: messageHandle,
	})
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Deleted")
	return nil
}

func main() {
	// Init Logger
	logger.InitLogger()
	logger.Agent.Info("Running Application BotFood")

	// Init Main DB
	db, err := database.ConnectMainDB()
	if err != nil {
		logger.Agent.Info(err.Error())
	}
	mainDB := database.NewDB(db)

	var consumerService service.ConsumerService

	consumerService = service.RegisterConsumerService(mainDB)

	timeout := flag.Int64("t", 5, "How long, in seconds, that the message is hidden from others")
	flag.Parse()

	if *timeout < 0 {
		*timeout = 0
	}

	if *timeout > 12*60*60 {
		*timeout = 12 * 60 * 60
	}
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(utils.GetEnv("AWS_SQS_REGION", "ap-southeast-1")),
		Credentials: credentials.NewStaticCredentials(utils.GetEnv("AWS_SQS_ACCESS_KEY_ID", "access_key_id"), utils.GetEnv("AWS_SQS_SECRET_ACCESS_KEY", "secret_access_key"), ""),
	})
	// snippet-end:[sqs.go.receive_messages.sess]

	// Get URL of queue
	urlResult, err := GetQueueURL(sess, utils.GetEnv("AWS_SQS_QUEUE_NAME", "queue-name"))
	if err != nil {
		fmt.Println("Got an error getting the queue URL:")
		fmt.Println(err)
		return
	}

	for true {
		queueURL := urlResult.QueueUrl
		msgResult, err := GetMessages(sess, queueURL, timeout)
		if err != nil {
			fmt.Println("Got an error receiving messages:")
			fmt.Println(err)
			return
		}

		if len(msgResult.Messages) > 0 {
			fmt.Println("Incoming Message ID:     " + *msgResult.Messages[0].MessageId)
			fmt.Println("Incoming Message ID:     " + *msgResult.Messages[0].Body)

			firstMessage := msgResult.Messages[0]
			messageBody := *firstMessage.Body
			if utils.IsJSON(messageBody) {
				var payload dto.ConsumerRequestPayload
				err = json.Unmarshal([]byte(messageBody), &payload)
				if err != nil {
					fmt.Println(err)
				}

				if payload.Type == "activity-outlet" {
					for key, data := range payload.DataOutlet {
						payload.DataOutlet[key].IsOpen = utils.Itob(data.PayloadIsOpen)
					}
				} else if payload.Type == "activity-item" {
					for key, data := range payload.DataItem.ItemNew {
						payload.DataItem.ItemNew[key].InStock = utils.Itob(data.PayloadInStock)
					}
					for key, data := range payload.DataItem.ItemChange {
						payload.DataItem.ItemChange[key].InStock = utils.Itob(data.PayloadInStock)
					}
					for key, data := range payload.DataItem.ItemDelete {
						payload.DataItem.ItemDelete[key].InStock = utils.Itob(data.PayloadInStock)
					}
				}
				payload.RawMessage = messageBody
				payload.MessageID = *firstMessage.MessageId
				payload.MessageSlug = *firstMessage.ReceiptHandle
				action := consumerService.ConsumeActivityMessage(payload)
				if action.Errors.HasErrors() {
					fmt.Println("Error!!")
					fmt.Println(action.Errors)
				} else {
					fmt.Println("SUCCESS")
				}
				DeleteMessage(sess, queueURL, msgResult.Messages[0].ReceiptHandle)
			} else {
				fmt.Println("Unknown message")
				fmt.Println(firstMessage.Body)
				DeleteMessage(sess, queueURL, msgResult.Messages[0].ReceiptHandle)
			}
		}

	}
}
