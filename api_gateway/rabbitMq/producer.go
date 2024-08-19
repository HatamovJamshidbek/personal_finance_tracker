package rabbitMq

import (
	"fmt"
	"github.com/streadway/amqp"
	_ "github.com/streadway/amqp"
)

type RabbitMqProduce interface {
	ProduceMassage(queueName string, message []byte) error
	Close()
}

type RabbitMqProducerInt struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMqProducerInt(url string) (*RabbitMqProducerInt, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	channel, err := conn.Channel()
	return &RabbitMqProducerInt{
		conn:    conn,
		channel: channel,
	}, nil
}
func (rabbitMq *RabbitMqProducerInt) ProduceMassage(queueName string, message []byte) error {
	declareQueue, err := rabbitMq.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return err
	}

	err = rabbitMq.channel.Publish(
		"",
		declareQueue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	if err != nil {
		return err
	}
	fmt.Printf(" [x] Sent", string(message))
	return nil
}
