package main

import (
	"time"

	route "github.com/ZooLearn/file/api/route"
	"github.com/ZooLearn/file/bootstrap"
	"github.com/ZooLearn/file/internal/log"
	"github.com/gin-gonic/gin"
)

func main() {

	app := bootstrap.App()
	defer func() {
		if err := app.Consumer.Shutdown(); err != nil {
			log.Errorf("shutdown consumer: %s", err)
		}
		if err := app.Producer.Shutdown(); err != nil {
			log.Errorf("shutdown producer: %s", err)
		}
	}()
	env := app.Env

	timeout := time.Duration(env.ContextTimeout) * time.Second

	gin := gin.Default()

	route.Setup(env, timeout, gin, app.Producer)
	if err := gin.Run(env.ServerAddress); err != nil {
		panic(err)
	}
}
