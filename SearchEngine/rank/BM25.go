package rank

import (
	"SearchEngine/utils"
	"fmt"
	movingaverage "github.com/RobinUS2/golang-moving-average"
	"math"
)

//BM25
//             IDF * TF * (k1 + 1)
//BM25 = sum ----------------------------
//            TF + k1 * (1 - b + b * D / L)
//
//                   总文档数目
//IDF = log2( ------------------------  + 1 )
//                出现该关键词的文档数目


//平均词数
type BM25 struct {
	 Parameters *BM25Parameter
	 Total	   uint32
	 //出现关键词的文档数目 用mapid计算
	 Ma       *movingaverage.MovingAverage
}

type BM25Parameter struct {
	K1 float32    //默认为2.0
	B  float32	  //默认为0.75
}

//按照关键词顺序
func (bm25 *BM25) GetIDF(frequence []uint32) []float32 {

	 var IDF []float32

	 for _, value := range frequence {
		 tmpvalue := float64(bm25.Total/value +1)
		 IDF = append(IDF,float32(math.Log2(tmpvalue)))
	 }
	 return IDF
}

func (bm25 *BM25) GetScore(keywords []string,words map[string] []int,IDF []float32) float32 {
	  score := float32(0)

	  for index, value := range keywords {
		  TF := float32(len(words[value]))
		  D := float32(len(words))
		  numerator := IDF[index]*TF*(bm25.Parameters.K1+1)
		  denominator := TF + bm25.Parameters.K1*(1-bm25.Parameters.B+bm25.Parameters.B*D/float32(bm25.Ma.Avg()))

		  score += numerator/denominator

	  }
	  return score

}

func (bm25 *BM25) WriteBM25(filepath string){

	 parameters := []float32{bm25.Parameters.K1,bm25.Parameters.B}
	 //存储 K1 B
	 utils.Write(&parameters,fmt.Sprintf("%s%d",filepath,1))

	 //存储 Total
	 utils.Write(&bm25.Total,fmt.Sprintf("%s%d",filepath,2))

	 //存储[]float64
	 nums := bm25.Ma.GetValues()

	 utils.Write(&nums,fmt.Sprintf("%s%d",filepath,3))
}


//返回 [K1 B] Total []float64
func (bm25 *BM25) ReadBM25(filepath string)([]float32,uint32,[]float64){

	KB := make([]float32,0)
	utils.Read(&KB, fmt.Sprintf("%s%d",filepath,1))

	var total  uint32
	utils.Read(&total, fmt.Sprintf("%s%d",filepath,2))

	nums := make([]float64,0)
	utils.Read(&nums, fmt.Sprintf("%s%d",filepath,3))
	return KB,total,nums
}