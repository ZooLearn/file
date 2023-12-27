package bootstrap

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ZooLearn/file/config"
	"github.com/ZooLearn/file/internal/log"
	"github.com/ZooLearn/file/internal/rabbitx"
	"github.com/ZooLearn/file/internal/transcodex"
	"github.com/ZooLearn/file/internal/tusdx"

	"github.com/rabbitmq/amqp091-go"
)

type Application struct {
	Env      config.EnvConf
	Producer *rabbitx.Producer
	Consumer *rabbitx.Consumer
}

func App() Application {
	cfgs, err := config.NewConfig("./config.yaml")
	if err != nil {
		log.Panicf("config.NewConfig error: %v", err)
	}

	producer := rabbitx.NewProducer(context.Background(), cfgs.ProducerConf)
	consumer := rabbitx.NewConsumer(context.Background(), cfgs.ConsumerConf, func(deliveries <-chan amqp091.Delivery, done chan error) {
		for {
			val := <-deliveries
			data := tusdx.Data{}
			if err := json.Unmarshal(val.Body, &data); err != nil {
				log.Errorf("unmarshal: %s", err)
				if err := val.Ack(false); err != nil {
					log.Infof("ack %s", err)
				}
				continue
			}
			if err := transcodex.CreateHLS(fmt.Sprintf("./uploads/%s", data.ID), data.ID, "./convert", 10); err != nil {
				panic(err)
			}
			if err := val.Ack(false); err != nil {
				log.Infof("ack %s", err)
			}
		}
	})
	app := &Application{
		Env:      cfgs.EnvConf,
		Producer: producer,
		Consumer: consumer,
	}

	return *app
}
