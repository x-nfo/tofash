package config

import (
	"fmt"

	"github.com/streadway/amqp"
)

func NewRabbitMQ(cfg RabbitMQ) (*amqp.Connection, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.User, cfg.Password, cfg.Host, cfg.Port)
	conn, err := amqp.Dial(url)
	if err != nil {
		fmt.Printf("[NewRabbitMQ] Failed to connect to RabbitMQ: %v", err)
		return nil, err
	}
	return conn, nil
}
