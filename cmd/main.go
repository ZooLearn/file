package main

import (
	"time"

	route "github.com/ZooLearn/file/api/route"
	"github.com/ZooLearn/file/bootstrap"
	"github.com/gin-gonic/gin"
)

func main() {

	app := bootstrap.App()
	defer func() {
		app.Consumer.Shutdown()
		app.Producer.Shutdown()
	}()
	env := app.Env

	timeout := time.Duration(env.ContextTimeout) * time.Second

	gin := gin.Default()

	route.Setup(env, timeout, gin, app.Producer)
	if err := gin.Run(env.ServerAddress); err != nil {
		panic(err)
	}
}
