package  main

import (
	"SearchEngine/core"
	"SearchEngine/router"
)

func main(){

	var pe core.PictureEngine

	pe.Init()

	router.Router(&pe)
}