package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/Abdulsametileri/package-tracking-app/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	QueueName = "package_status"
)

type rabbitmqClient struct {
	conn          *amqp.Connection
	ch            *amqp.Channel
	connString    string
	packageStatus <-chan amqp.Delivery
}

func NewRabbitMQClient(connectionString string) (*rabbitmqClient, error) {
	c := &rabbitmqClient{}
	var err error

	c.conn, err = amqp.Dial(connectionString)
	if err != nil {
		return nil, err
	}

	c.ch, err = c.conn.Channel()
	if err != nil {
		return nil, err
	}

	err = c.configureQueue()

	return c, err
}

func (c *rabbitmqClient) ConsumeByVehicleID(ctx context.Context, vehicleID string) ([]byte, error) {
	for msg := range c.packageStatus {
		if msg.MessageId == vehicleID {
			return msg.Body, nil
		}
	}
	return nil, errors.New("err when getting package status on channel")
}

func (c *rabbitmqClient) Publish(p domain.Package) {
	jsonStr := fmt.Sprintf(`{ "from": "%s", "to": "%s", "vehicleId": "%s" }`, p.From, p.To, p.VehicleID)

	c.ch.Publish("", QueueName, false, false, amqp.Publishing{
		ContentType: "application/json",
		MessageId:   p.VehicleID,
		Body:        []byte(jsonStr),
	})
}

func (c *rabbitmqClient) Close() {
	c.ch.Close()
	c.conn.Close()
}

func (c *rabbitmqClient) configureQueue() error {
	_, err := c.ch.QueueDeclare(
		QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	c.packageStatus, err = c.ch.Consume(
		QueueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	return err
}
