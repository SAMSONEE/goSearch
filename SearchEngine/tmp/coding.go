package main

import (
	"SearchEngine/utils"
	"fmt"
)

type Vector struct {
	x, y, z int
}

func main() {
	v := Vector{
		x:1,
		y:2,
		z:3,
	}

	data := utils.Encoder(v)

	fmt.Println(data)

	var  result Vector
	utils.Decoder(data, &result)

	fmt.Println(result)
}