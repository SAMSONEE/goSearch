package core

import (
	"fmt"
	"sort"
	"sync"
)

//  排序 按照分数
//type Score struct {
//	Id 	uint32
//	Score float32
//}




type ScoreSlice []Score

func (s ScoreSlice) Len() int {
	return len(s)
}

func (s ScoreSlice) Less(i, j int) bool {
	return s[i].Score < s[j].Score
}

func (s ScoreSlice) Swap(i,j int) {
	s[i],s[j] = s[j], s[i]
}

type FastSort struct {
	data []*Score
	sync.Mutex
}

type FastType []*Score


func (ft FastType) Len() int {
	return len(ft)
}

func (ft FastType) Less(i, j int) bool {
	return ft[i].Id < ft[j].Id
}

func (ft FastType) Swap(i,j int) {
	ft[i],ft[j] = ft[j], ft[i]
}

func find(data []Score, id uint32) (bool, int) {
	left := 0
	right := len(data) - 1

	for left <= right {
		mid := left + (right-left)/2
		if data[mid].Id < id {
			left = mid+1
		}else if  data[mid].Id > id{
			right = mid -1
		}else {
			return true, mid
		}
	}
	return false, -1
}

func (f *FastSort) Count() int {
	return len(f.data)
}

func (f *FastSort) Add(scores []*Score) {
	 if scores == nil {
		 return
	 }
	 f.Lock()
	 defer f.Unlock()
	 f.data = append(f.data,scores...)
}

//不用二分查找 因为f.data 已经排序
func (f *FastSort) GetSort()  []Score {
	 var result = make([]Score,len(f.data))

	 //协程
	 sort.Sort(FastType(f.data))


	 k := 0

	 //fmt.Println(f.data)

	 for _, item := range f.data{
		 flag, index := find(result,item.Id)
		 if flag == true {
			 result[index].Score += item.Score
		 }else {
			 result[k] = Score{
				 Id: item.Id,
				 Score: item.Score,
			 }
		 }
		 k++
	 }

	 fmt.Println("*************")

	 sort.Sort(sort.Reverse(ScoreSlice(result)))

	 //for _, item :=  range result {
		//fmt.Println("Id:",item.Id," Score:",item.Score)
	 //}

	 return result[:k]

}















































