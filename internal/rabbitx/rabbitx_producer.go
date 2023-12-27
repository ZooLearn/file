package rabbitx

import (
	"context"

	"github.com/ZooLearn/file/internal/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitmqProConf struct {
	Address      string `yaml:"Address"`
	Exchange     string `yaml:"Exchange"`
	ExchangeType string `yaml:"ExchangeType"`
	Queue        string `yaml:"Queue"`
	BindingKey   string `yaml:"BindingKey"`
	ConsumerTag  string `yaml:"ConsumerTag"`
	RoutingKey   string `yaml:"RoutingKey"`
	LifeTime     int    `yaml:"LifeTime"`
	AutoACK      bool   `yaml:"AutoACK"`
}

type Producer struct {
	ctx     context.Context
	cfgs    RabbitmqProConf
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewProducer(ctx context.Context, cfgs RabbitmqProConf) *Producer {
	config := amqp.Config{
		Vhost:      "/",
		Properties: amqp.NewConnectionProperties(),
	}
	config.Properties.SetClientConnectionName("producer_sv_1")
	log.Infof("producer: dialing %s", cfgs.Address)
	conn, err := amqp.DialConfig(cfgs.Address, config)
	if err != nil {
		log.Errorf("producer: error in dial: %s", err)
		panic(err)
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Errorf("error getting a channel: %s", err)
		panic(err)
	}

	log.Infof("producer: declaring exchange")
	if err := channel.ExchangeDeclare(
		cfgs.Exchange,     // name
		cfgs.ExchangeType, // type
		true,              // durable
		false,             // auto-delete
		false,             // internal
		false,             // noWait
		nil,               // arguments
	); err != nil {
		log.Errorf("producer: Exchange Declare: %s", err)
		panic(err)
	}
	log.Infof("producer: declaring queue '%s'", cfgs.Queue)
	queue, err := channel.QueueDeclare(
		cfgs.Queue, // name of the queue
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		log.Errorf("producer: Queue Declare: %s", err)
		panic(err)
	}
	log.Infof("producer: declared queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		queue.Name, queue.Messages, queue.Consumers, cfgs.RoutingKey)
	log.Infof("producer: declaring binding")
	if err := channel.QueueBind(queue.Name, cfgs.RoutingKey, cfgs.Exchange, false, nil); err != nil {
		log.Errorf("producer: Queue Bind: %s", err)
		panic(err)
	}
	return &Producer{
		ctx:     ctx,
		cfgs:    cfgs,
		conn:    conn,
		channel: channel,
	}
}

func (p *Producer) Publish(data []byte) error {

	if err := p.channel.PublishWithContext(context.Background(), p.cfgs.Exchange, p.cfgs.RoutingKey, true, false, amqp.Publishing{
		Body: data,
	}); err != nil {
		log.Errorf("send data to queue: %s", err)
		return err
	}
	return nil
}

func (p *Producer) Shutdown() error {
	if err := p.channel.Close(); err != nil {
		return err
	}
	if err := p.conn.Close(); err != nil {
		return err
	}
	return nil
}
