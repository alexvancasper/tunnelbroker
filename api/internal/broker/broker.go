package broker

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type MsgBroker struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func (m *MsgBroker) connect(connStr string) error {
	conn, err := amqp.Dial(connStr) //"amqp://guest:guest@localhost:5672/"
	if err != nil {
		return fmt.Errorf("%s: %s", "Failed to connect to RabbitMQ", err)
	}
	m.conn = conn
	return nil
}

func (m *MsgBroker) createChannel() error {
	ch, err := m.conn.Channel()
	if err != nil {
		return fmt.Errorf("%s: %s", "Failed to open a channel", err)
	}
	m.channel = ch
	return nil
}

func (m *MsgBroker) queueDeclare(queueName string) error {
	var err error
	m.queue, err = m.channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("%s: %s", "Failed to declare a queue", err)
	}
	return nil
}

func (m *MsgBroker) PublishMsg(data []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	err = m.channel.PublishWithContext(ctx,
		"",           // exchange
		m.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		})
	if err != nil {
		return fmt.Errorf("%s: %s", "Failed to publish a message", err)
	}
	log.Printf(" [x] Sent %s\n", data)
	return nil
}

func (m *MsgBroker) connClose() {
	_ = m.conn.Close()
}
func (m *MsgBroker) channelClose() {
	_ = m.channel.Close()
}

func MsgBrokerInit(connStr, queueName string) (*MsgBroker, error) {
	var msg MsgBroker
	var err error
	err = msg.connect(connStr)
	if err != nil {
		return nil, err
	}
	err = msg.createChannel()
	if err != nil {
		return nil, err
	}
	err = msg.queueDeclare(queueName)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (m *MsgBroker) Close() {
	m.channelClose()
	m.connClose()
}
