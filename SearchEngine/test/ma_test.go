package test

import (
	"SearchEngine/rank"
	"fmt"
	movingaverage "github.com/RobinUS2/golang-moving-average"
	"math/rand"
	"testing"
	"time"
)
func Ma(){

	ma := movingaverage.New(20)

	pa := &rank.BM25Parameter{
		K1: 2.0,
		B:  0.75,
	}

	bm25 := &rank.BM25{
		Parameters: pa,
		Total: uint32(0),
		Ma: ma,
	}

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 40; i++ {
		bm25.Ma.Add(rand.Float64())
	}

	filepath := "/home/wqk/GolandProjects/SearchEngine/test/testBM25/BM25"

	fmt.Println("K1:",bm25.Parameters.K1," B:",bm25.Parameters.B)
	fmt.Println("Total:",bm25.Total)
	fmt.Println("values:")
	fmt.Println(bm25.Ma.GetValues())

	bm25.WriteBM25(filepath)

	ma2 := movingaverage.New(20)

	var bm rank.BM25

	Pa, Total, values := bm.ReadBM25(filepath)

	pa2 := &rank.BM25Parameter{
		K1: Pa[0],
		B: Pa[1],
	}

	bm.Total = Total
	bm.Parameters = pa2
	ma2.Addnums(values)

	bm.Ma = ma

	fmt.Println("**************************")
	fmt.Println(bm.Parameters)
	fmt.Println("Total:",bm.Total)
	fmt.Println("values:")
	fmt.Println(bm.Ma.GetValues())


}


func TestMa(t *testing.T) {
	Ma()
}