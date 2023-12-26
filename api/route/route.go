package route

import (
	"time"

	"github.com/ZooLearn/file/internal/rabbitx"
	"github.com/ZooLearn/file/internal/tusdx"

	"github.com/ZooLearn/file/config"
	"github.com/gin-gonic/gin"
)

func Setup(env config.EnvConf, timeout time.Duration, ginx *gin.Engine, producer *rabbitx.Producer) {
	protectedRouter := ginx.Group("")
	protectedRouter.Use()
	handler := tusdx.TusdHandler(producer)
	protectedRouter.POST("/files/", gin.WrapF(handler.PostFile))
	protectedRouter.HEAD("/files/:id", gin.WrapF(handler.HeadFile))
	protectedRouter.PATCH("/files/:id", gin.WrapF(handler.PatchFile))
	protectedRouter.GET("/files/:id", gin.WrapF(handler.GetFile))
}
