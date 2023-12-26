package rabbitx

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ZooLearn/file/internal/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

/*
	uri               = flag.String("uri", "amqp://guest:guest@localhost:5672/", "AMQP URI")
	exchange          = flag.String("exchange", "test-exchange", "Durable, non-auto-deleted AMQP exchange name")
	exchangeType      = flag.String("exchange-type", "direct", "Exchange type - direct|fanout|topic|x-custom")
	queue             = flag.String("queue", "test-queue", "Ephemeral AMQP queue name")
	bindingKey        = flag.String("key", "test-key", "AMQP binding key")
	consumerTag       = flag.String("consumer-tag", "simple-consumer", "AMQP consumer tag (should not be blank)")
	lifetime          = flag.Duration("lifetime", 5*time.Second, "lifetime of process before shutdown (0s=infinite)")
	autoAck           = flag.Bool("auto_ack", false, "enable message auto-ack")
*/

type RabbitmqConConf struct {
	Address      string `yaml:"Address"`
	Exchange     string `yaml:"Exchange"`
	ExchangeType string `yaml:"ExchangeType"`
	Queue        string `yaml:"Queue"`
	BindingKey   string `yaml:"BindingKey"`
	ConsumerTag  string `yaml:"ConsumerTag"`
	LifeTime     int    `yaml:"LifeTime"`
	VerBose      bool   `yaml:"VerBose"`
	AutoACK      bool   `yaml:"AutoACK"`
}

type Consumer struct {
	ctx     context.Context
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
	pull    func(deliveries <-chan amqp.Delivery, done chan error)
}

func SetupCloseHandler(consumer *Consumer) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Info("Ctrl+C pressed in Terminal")
		if err := consumer.Shutdown(); err != nil {
			log.Errorf("error during shutdown: %s", err)
		}
		os.Exit(0)
	}()
}

func NewConsumer(ctx context.Context, cfgs RabbitmqConConf, handle func(deliveries <-chan amqp.Delivery, done chan error)) *Consumer {
	c := &Consumer{
		ctx:     ctx,
		conn:    nil,
		channel: nil,
		pull:    handle,
		done:    make(chan error),
	}

	var err error

	config := amqp.Config{Properties: amqp.NewConnectionProperties()}
	config.Properties.SetClientConnectionName("consumer_sv_1")
	log.Infof("dialing %q", cfgs.Address)
	c.conn, err = amqp.DialConfig(cfgs.Address, config)
	if err != nil {
		panic(err)
		// return nil, fmt.Errorf("dial: %s", err)
	}

	go func() {
		log.Infof("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
	}()

	log.Info("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		log.Errorf("hannel: %s", err)
		panic(err)
	}

	log.Infof("got Channel, declaring Exchange (%q)", cfgs.Exchange)
	if err = c.channel.ExchangeDeclare(
		cfgs.Exchange,     // name of the exchange
		cfgs.ExchangeType, // type
		true,              // durable
		false,             // delete when complete
		false,             // internal
		false,             // noWait
		nil,               // arguments
	); err != nil {
		log.Errorf("exchange declare: %s", err)
		panic(err)
	}

	log.Infof("declared Exchange, declaring Queue %q", cfgs.Queue)
	queue, err := c.channel.QueueDeclare(
		cfgs.Queue, // name of the queue
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		log.Errorf("queue declare: %s", err)
		panic(err)
	}

	log.Infof("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		queue.Name, queue.Messages, queue.Consumers, cfgs.BindingKey)

	if err = c.channel.QueueBind(
		queue.Name,      // name of the queue
		cfgs.BindingKey, // bindingKey
		cfgs.Exchange,   // sourceExchange
		false,           // noWait
		nil,             // arguments
	); err != nil {
		log.Errorf("queue bind: %s", err)
		panic(err)
	}

	log.Infof("Queue bound to Exchange, starting Consume (consumer tag %q)", c.tag)
	deliveries, err := c.channel.Consume(
		queue.Name,   // name
		c.tag,        // consumerTag,
		cfgs.AutoACK, // autoAck
		false,        // exclusive
		false,        // noLocal
		false,        // noWait
		nil,          // arguments
	)
	if err != nil {
		log.Errorf("queue consume: %s", err)
		panic(err)
	}

	if c.pull == nil {
		log.Error("missing handle function")
		panic(err)
	}
	go c.pull(deliveries, c.done)

	return c
}

func (c *Consumer) Shutdown() error {
	// will close() the deliveries channel
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer log.Infof("AMQP shutdown OK")

	// wait for handle() to exit
	return <-c.done
}
