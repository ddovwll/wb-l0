package message_querry

import (
	"context"
	"demoService/src/application/services"
	"demoService/src/domain/models"
	"encoding/json"
	"log"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/segmentio/kafka-go"
)

type OrderConsumer struct {
	orderService services.OrderService
	validate     *validator.Validate
	reader       *kafka.Reader
	dlqWriter    *kafka.Writer
	wg           sync.WaitGroup
}

func NewOrderConsumer(orderService services.OrderService, reader *kafka.Reader, dlqWriter *kafka.Writer) *OrderConsumer {
	return &OrderConsumer{
		orderService: orderService,
		validate:     validator.New(),
		reader:       reader,
		dlqWriter:    dlqWriter,
	}
}

func (orderConsumer *OrderConsumer) StartConsumer(context context.Context, workerCount int) {
	jobs := make(chan kafka.Message, 100)
	defer close(jobs)

	for i := 0; i < workerCount; i++ {
		orderConsumer.wg.Add(1)
		go func() {
			defer orderConsumer.wg.Done()
			for msg := range jobs {
				orderConsumer.processOrder(context, msg)
			}
		}()
	}

	log.Println("Kafka consumer started")

	for {
		msg, err := orderConsumer.reader.FetchMessage(context)
		if err != nil {
			if context.Err() != nil {
				log.Println("Consumer stopped")
				break
			}
			log.Println("Error reading Kafka message:", err)
			continue
		}
		jobs <- msg
	}

	orderConsumer.wg.Wait()

	if err := orderConsumer.reader.Close(); err != nil {
		log.Println("Error closing Kafka reader:", err)
	}
	log.Println("Consumer shutdown complete")
}

func (orderConsumer *OrderConsumer) processOrder(context context.Context, message kafka.Message) {
	var order models.Order
	if err := json.Unmarshal(message.Value, &order); err != nil {
		log.Println("Failed to parse order", err)
		orderConsumer.sendToDLQ(context, message, err)
		if commitErr := orderConsumer.reader.CommitMessages(context, message); commitErr != nil {
			log.Println("Commit error after DLQ:", commitErr)
		}
		return
	}

	if err := orderConsumer.validate.Struct(order); err != nil {
		log.Printf("Validation error order %s: %s", order.OrderUID, err)
		orderConsumer.sendToDLQ(context, message, err)
		if commitErr := orderConsumer.reader.CommitMessages(context, message); commitErr != nil {
			log.Println("Commit error after DLQ:", commitErr)
		}
		return
	}

	log.Println("Consumed order with id:", order.OrderUID)

	if err := orderConsumer.orderService.Create(context, order); err != nil {
		log.Println("Failed to save order:", err)
		orderConsumer.sendToDLQ(context, message, err)
		if commitErr := orderConsumer.reader.CommitMessages(context, message); commitErr != nil {
			log.Println("Commit error after DLQ:", commitErr)
		}
		return
	}

	if err := orderConsumer.reader.CommitMessages(context, message); err != nil {
		log.Println("Commit error:", err)
	}
}

func (orderConsumer *OrderConsumer) sendToDLQ(ctx context.Context, message kafka.Message, error error) {
	dlqMessage := kafka.Message{
		Value: message.Value,
		Headers: []kafka.Header{
			{
				Key:   "X-Error",
				Value: []byte(error.Error()),
			},
		},
	}

	if err := orderConsumer.dlqWriter.WriteMessages(ctx, dlqMessage); err != nil {
		log.Println("Failed to write to DLQ:", err)
	} else {
		log.Println("Message sent to DLQ")
	}
}
