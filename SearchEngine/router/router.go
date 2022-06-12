package router

import (
	"SearchEngine/core"
	"SearchEngine/router/api"
	"github.com/gin-gonic/gin"
)

func Router(pe *core.PictureEngine){

	router := gin.Default()

	api.Register(router)

	api.Set(pe)


	router.Run(":8080")

}
