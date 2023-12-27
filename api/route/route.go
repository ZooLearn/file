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
	handler := tusdx.TusdMediaHandler(producer)
	protectedRouter.POST("/media/files/", gin.WrapF(handler.PostFile))
	protectedRouter.HEAD("/media/files/:id", gin.WrapF(handler.HeadFile))
	protectedRouter.PATCH("/media/files/:id", gin.WrapF(handler.PatchFile))
	protectedRouter.GET("/media/files/:id", gin.WrapF(handler.GetFile))
}
