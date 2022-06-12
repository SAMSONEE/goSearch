package test

import (
	"SearchEngine/core"
	"fmt"
	"testing"
)


func Init2() *core.PictureEngine{
	var pe core.PictureEngine

	pe.Init()

	//filepath := "/home/wqk/GolandProjects/SearchEngine/test/data/wukong50k_release.csv"
	//
	//
	//pe.SegmentCsv(filepath)

	//字典树写到磁盘
	//fileTrie := "/home/wqk/GolandProjects/SearchEngine/test/data/storetrie/st"
	//trie.Write(pe.Tire,fileTrie)

	return &pe

}

func Search(){

	pe := Init2()

	//fmt.Println(pe.KeyMapId)
	//fmt.Println(pe.IdMapKey)
	//fmt.Println(pe.IdMapDocument)



	 //filepath := "/home/wqk/GolandProjects/SearchEngine/test/data/wukong50k_release.csv"
	 //
	 //pe.SegmentCsv(filepath)

	 request := core.Searchrequest{
		QueryText: "红色我爱你，我的家乡",
		KeyWords: "天空",
		HateWords: "红色",
	 }

	 result :=pe.Search(request)

	 //fmt.Println(len(result))
	 for _ ,item := range result {
		 fmt.Println(len(item.Pictures))
	 }
}

func TestSearch(t *testing.T){
	Search()
}