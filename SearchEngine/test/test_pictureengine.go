package main

import (
	"SearchEngine/core"
	"SearchEngine/utils"
	"fmt"
)

//func main() {
//	var segmenter  sego.Segmenter
//	segmenter.LoadDictionary("/home/wqk/GolandProjects/SearchEngine/test/data/dictionary.txt")
//}


//1 分词

//func main(){
//	 var pEngine PictureEngine
//
//	 pEngine.Init()
//	 //
//	 //filepath := "/home/wqk/GolandProjects/SearchEngine/test/data/wukong50k_release.csv"
//	 //
//	 //pEngine.SegmentCsv(filepath)
//
//}


//2 字典树索引

//func main(){
//
//	tree := NewTrie()
//	tree.InsertWord("我和我的祖国")
//	tree.InsertWord("我爱我的祖国")
//	tree.InsertWord("我和我的家乡")
//	fmt.Println(tree.Contains("我"))
//	fmt.Println(tree.Contains("我爱"))
//	fmt.Println(tree.Contains("我和"))
//	fmt.Println(tree.Contains("我和我的家乡"))
//	//fmt.Println(tree.IsPrefix("我"))
//	//fmt.Println(tree.IsPrefix("我爱"))
//	//fmt.Println(tree.IsPrefix("我和"))
//	fmt.Println(tree.GetSize())
//}


//func main() {
//
//	a := 40506
//	b := 23434
//
//	s1 := string(a+b)
//	s2 := string(a)+string(b)
//
//	fmt.Println(s1)
//	fmt.Println(s2)
//}

//3 leveldb

//func main () {
//	filepath := "SearchEngine/data/db"
//
//	db, _ := leveldb.OpenFile(filepath,nil)
//
//	defer db.Close()
//
//	db.Put([]byte("key"), []byte("value"), nil)
//
//	key, _ := db.Get([]byte("key"),nil)
//
//	fmt.Println(key)
//
//}



// 4 coding

//type Vector struct {
//	X, Y, Z int
//}
//
//func main() {
//	v := Vector{
//		X:1,
//		Y:2,
//		Z:3,
//	}
//
//	data := utils.Encoder(v)
//
//	fmt.Println(data)
//
//	var  result Vector
//	utils.Decoder(data, &result)
//
//	fmt.Println(result.Z)
//}


// 5 string to int

//func main() {
//
//	a := "我和你"
//
//	result := utils.StringToInt(a)
//
//	fmt.Println(result)
//
//
//}

//leveldb

//func main(){
//
//	path := "SearchEngine/data/db"
//
//	ldb, err := leveldb.CreateLeveldb(path)
//
//	if err != nil {
//		return
//	}
//
//	key1 := "nihao"
//	//key2 := "nihao2"
//	//val1 := "231"
//	//val2 := "123"
//	//
//	key1byte := utils.Uint32ToBytes(utils.StringToInt(key1))
//	//key2byte := utils.Uint32ToBytes(utils.StringToInt(key2))
//	//
//	//val1byte := utils.Uint32ToBytes(utils.StringToInt(val1))
//	//fmt.Println(utils.StringToInt(val1))
//	//val2byte := utils.Uint32ToBytes(utils.StringToInt(val2))
//	//
//	//err = ldb.Put(key1byte,val1byte)
//	//if err != nil {
//	//	return
//	//}
//	//
//	//err = ldb.Put(key2byte,val2byte)
//	//
//	//if err != nil {
//	//	return
//	//}
//
//	val, err := ldb.Get(key1byte)
//
//	if err != nil {
//		return
//	}
//
//	fmt.Println(utils.BytesToUint32(val))
//
//	flag, err := ldb.Has([]byte("12323"))
//
//	if err != nil {
//		return
//	}
//
//	fmt.Println(flag)
//
//	fmt.Println(utils.BytesToUint32(val))
//
//	flag2, err := ldb.Has(key1byte )
//
//	if err != nil {
//		return
//	}
//
//	fmt.Println(flag2)
//
//
//	err =ldb.Close()
//
//	if err != nil {
//		return
//	}
//}


//tire


//func main (){
//	var KeyMapId []*leveldb.Leveldb
//
//	Keypath := "SearchEngine/data/db/KeyMapId/"
//
//	for i := 0 ; i < 10; i++ {
//		KeyMapId = append(KeyMapId,leveldb.CreateLeveldb(Keypath+string(i)))
//		println(i)
//	}
//}


func main() {
	var pe core.PictureEngine

	pe.Init()

	filepath := "/home/wqk/GolandProjects/SearchEngine/test/data/wukong50k_release.csv"

	pe.SegmentCsv(filepath)

	keyword := "短袖"
	keyint := utils.StringToInt(keyword)
	k := utils.Uint32ToBytes(keyint)

	ldb := pe.KeyMapId[pe.GetLeveldbId(keyint)]

	fmt.Println(pe.Tire.Contains(keyword))

	flag, _ := ldb.Has(k)

	fmt.Println(flag)

	if flag {
		ids := make([]uint32,0)
		buf, _ := ldb.Get(k)
		utils.Decoder(buf,&ids)
		fmt.Println(ids)
	}
}


//func main() {
//	Keypath := "/home/wqk/GolandProjects/SearchEngine/test/data/db/KeyMapId/"
//	Idpath := "/home/wqk/GolandProjects/SearchEngine/test/data/db/IdMapKey/"
//	Docpath := "/home/wqk/GolandProjects/SearchEngine/test/data/db/IdMapDocument/"
//
//
//	//初始化leveldb
//	for i := 0 ; i < 10; i++ {
//		fmt.Println(fmt.Sprintf("%s%d",Keypath,i))
//		fmt.Println(fmt.Sprintf("%s%d",Idpath,i))
//		fmt.Println(fmt.Sprintf("%s%d",Docpath,i))
//	}
//}


// 序列化

//func main () {
//
//	tree := NewTrie()
//
//	word := "未全开"
//	word1 := "未来"
//	word2 := "未全亿"
//	word3 := "未来来来"
//	tree.InsertWord(word)
//	tree.InsertWord(word1)
//	tree.InsertWord(word2)
//	tree.InsertWord(word3)
//
//	fmt.Println(tree.Level)
//
//	filename := "/home/wqk/GolandProjects/SearchEngine/test/data/storetrie/st"
//
//	Write(tree,filename)
//
//	newtree := Read(filename)
//
//	Write(newtree,filename)
//}


//func main () {
//	filename := "/home/wqk/GolandProjects/SearchEngine/test/data/storetrie/st"
//	err := ioutil.WriteFile(filename,[]byte("123"),0600)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//
//}


//fastsort


//func main(){
//
//
//
//}
//




































































