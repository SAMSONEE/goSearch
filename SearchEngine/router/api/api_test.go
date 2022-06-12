package api

import (
	"SearchEngine/core"
	"github.com/gin-gonic/gin"
	"testing"
)


func Init(){

	var pe core.PictureEngine

	router := gin.Default()

	Register(router)

	pe.Init()

	Set(&pe)


	router.Run(":8080")
}



func TestInit(t *testing.T) {
	Init()
}