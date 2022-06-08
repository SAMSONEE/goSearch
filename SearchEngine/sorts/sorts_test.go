package sorts

import (
	"SearchEngine/core"
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

func add(count uint64) (uint64,*FastSort){
	rand.Seed(time.Now().UnixNano())

	fastsort := new(FastSort)

	for  i := 0 ;i < math.MaxUint16; i++ {

		size := rand.Intn(5)

		scores := make([]*core.Score,0)
		for j := 0; j < size ;j++ {
			//计数器
			count++

			randid := rand.Uint32()
			randscore := rand.Float32()
			score := &core.Score{
				Id: randid,
				Score: randscore,
			}
			scores = append(scores,score)
		}
		fastsort.Add(scores)
	}
	return count,fastsort
}


func TestAdd(t *testing.T) {

	 count, fs := add(uint64(0))

	 fmt.Println(count)

	 result := fs.GetSort()

	 for index, sc := range result {
		if index > 10 {
			 break
		}
		fmt.Println("Id:",sc.Id," Score:",sc.Score)
	 }
}