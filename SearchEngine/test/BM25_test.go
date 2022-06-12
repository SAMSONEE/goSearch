package test

import (
	"SearchEngine/core"
	"SearchEngine/utils"
	"fmt"
	"testing"
)

func Init() *core.PictureEngine{
		var pe core.PictureEngine

		pe.Init()

		filepath := "/home/wqk/GolandProjects/SearchEngine/test/data/wukong50k_release.csv"


		pe.SegmentCsv(filepath)

		return &pe

}

func TestInit(t *testing.T) {

	pe := Init()

	keywords := []string {"红色","蓝天"}

	ids := make([]uint32,0)

	frequence := make([]uint32,0)

	for _, word := range keywords {
		keyint := utils.StringToInt(word)
		k := utils.Uint32ToBytes(keyint)
		ldb := pe.KeyMapId[pe.GetLeveldbId(keyint)]

		if pe.Tire.Contains(word){
			flag, _ := ldb.Has(k)
			tmpid := make([]uint32,0)
			if flag {
				buf, _ := ldb.Get(k)
				utils.Decoder(buf,&tmpid)
			}
			frequence = append(frequence,uint32(len(tmpid)))
			ids = append(ids,tmpid...)
		}
	}

	//fmt.Println("Ids:")
	//fmt.Println(ids)
	//
	//fmt.Println("********************")

	fastsort := new(core.FastSort)
	Scs := make([]*core.Score,0)

	for _, id := range ids {

		ldb2 := pe.IdMapKey[pe.GetLeveldbId(id)]

		k := utils.Uint32ToBytes(id)

		flag, _ := ldb2.Has(k)

		if flag {
			keymap := make(map[string][]int,0)
			buf, _:= ldb2.Get(k)
			utils.Decoder(buf,&keymap)
			score := pe.GetRank(keywords,keymap,frequence)
			//fmt.Println(id ," ",score)
			Sc := &core.Score{
				Id: id,
				Score: score,
			}
			Scs = append(Scs,Sc)
		}
	}
	fastsort.Add(Scs)
	//fmt.Println("******************")
	result := fastsort.GetSort()

	for _, value := range result {
		fmt.Println("Id:",value.Id," Score:",value.Score)
	}

}
