/*
*   Todo:
*        1. 得分
*        2. 初始化
*        3. 文档存储格式（id：JSON格式）
*
 */

package core

import (
	"SearchEngine/leveldb"
	"SearchEngine/rank"
	. "SearchEngine/trie"
	"SearchEngine/utils"
	"fmt"
	movingaverage "github.com/RobinUS2/golang-moving-average"
	"github.com/huichen/sego"
	"log"
	"sync"
)



type PictureEngine struct {

	//被索引的文档数
	numDocumentIsIndex uint32
	//被分词的文档数
	numSegment uint32
	//被存储的文档数
	numDocumentStored uint32

	//被存储的文档的平均词数
	numAvergeword uint8

	//停字符
	stoptokens StopTokens
	//分词器
	segmenter sego.Segmenter

	//BM25
	bm25 *rank.BM25


	//字典树索引
	Tire *TrieTree

	//leveldb 分十组
	KeyMapId		[]*leveldb.Leveldb
	IdMapKey		[]*leveldb.Leveldb
	IdMapDocument   []*leveldb.Leveldb

	//锁
	sync.Mutex
	//等待
	sync.WaitGroup
}

//初始化
func (pe *PictureEngine) Init() {
	//加载词典
	pe.segmenter.LoadDictionary("/home/wqk/GolandProjects/SearchEngine/test/data/dictionary.txt")

	//通用词路径
	stopname := "/home/wqk/GolandProjects/SearchEngine/test/data/stop_tokens.txt"
	//初始化停用词
	pe.stoptokens.ReadStop(stopname)

	//绝对路径
	Keypath := "/home/wqk/GolandProjects/SearchEngine/test/data/db/KeyMapId/"
	Idpath := "/home/wqk/GolandProjects/SearchEngine/test/data/db/IdMapKey/"
	Docpath := "/home/wqk/GolandProjects/SearchEngine/test/data/db/IdMapDocument/"


	//初始化leveldb
	for i := 0 ; i < 10; i++ {
		pe.KeyMapId = append(pe.KeyMapId,leveldb.CreateLeveldb(fmt.Sprintf("%s%d",Keypath,i)))
		pe.IdMapKey = append(pe.IdMapKey,leveldb.CreateLeveldb(fmt.Sprintf("%s%d",Idpath,i)))
		pe.IdMapDocument = append(pe.IdMapDocument,leveldb.CreateLeveldb(fmt.Sprintf("%s%d",Docpath,i)))
	}

	//初始化字典树
	pe.Tire = NewTrie()

	//for i:= 0; i < 10; i++ {
	//	ldb := pe.KeyMapId[i]
	//	iter := ldb.db.NewIterator(nil,nil)
	//
	//	for iter.Next() {
	//		keyword := utils.BytesToUint32(iter.Key())
	//	}
	//}

	//初始化BM25 有点问题
	//ma := movingaverage.New(math.MaxUint32)
	//pa := &rank.BM25Parameter{
	//	K1: 2.0,
	//	B:  0.75,
	//}
	//pe.bm25 =  &rank.BM25{
	//	Parameters: pa,
	//	Total: 0,
	//	Ma: ma,
	//}

	ma := movingaverage.New(60)
	pa := &rank.BM25Parameter{
			K1: 2.0,
			B:  0.75,
	}

	pe.bm25 = &rank.BM25{
		Parameters: pa,
		Total: uint32(0),
		Ma: ma,
	}

}

func (pe *PictureEngine) SegmentCsv(filePath string){

	//读取CSV文件
	message, err := Readcsv(filePath)

	if err != nil {
		log.Fatalln(err)
	}

	//如何将string 转化为byte
	// t[0] url t[1] text
	for _, t := range message {

		pe.numSegment++

		tokensmap := make(map[string] []int)

		text := []byte(t[1])
		segments := pe.segmenter.Segment(text)
		// words := sego.SegmentsToSlice(segments, true)
		for _, segment := range segments {
			token := segment.Token().Text()

			if !pe.stoptokens.Stop_Tokens[token] {
				tokensmap[token] = append(tokensmap[token], segment.Start())
			}
		}
		//添加倒排索引 key->id
		pe.numDocumentIsIndex++

		//BM25
		pe.bm25.Total++
		pe.bm25.Ma.Add(float64(len(tokensmap)))



		for word, _ := range tokensmap {
			pe.AddIndexOfKeyId(word,pe.numDocumentIsIndex)

			////Debug
			//fmt.Print(word," ")
		}
		//添加正排索引
		pe.AddIndexOfIdKey(pe.numDocumentIsIndex,tokensmap)
		//添加文档
		picture := &Picture{
			Id:pe.numDocumentIsIndex,
			Picture_url: t[0],
			Picture_context: t[1],
		}
		pe.AddDocument(pe.numDocumentIsIndex,picture)

		fmt.Println(pe.numSegment,", id:",pe.numDocumentIsIndex)

	}
}

func (pe *PictureEngine) GetLeveldbId(id uint32) int {
	return int(id% 10)
}


//添加倒排索引
/*
*   1.判断是否在字典树
*	2.判断是否在索引里
*/

// 序列化后的keyword ：id (uint32)

func (pe *PictureEngine) AddIndexOfKeyId(keyword string, id uint32){
	pe.Lock()

	defer pe.Unlock()

	keyint := utils.StringToInt(keyword)
	k := utils.Uint32ToBytes(keyint)
	ldb := pe.KeyMapId[pe.GetLeveldbId(keyint)]
	ids := make([]uint32,0)

	//是否存在字典树
	//不在 放入
	if pe.Tire.Contains(keyword) {
		flag, _ :=ldb.Has(k)
		if flag {
			buf, _ := ldb.Get(k)
			utils.Decoder(buf,&ids)

			//判断是否存在
			exit := false
			for _, v := range ids {
				if v==id {
					exit = true
					break
				}
			}

			if !exit {
				ids = append(ids,id)
			}
		}else {
			ids = append(ids,id)
		}
	}else {
		pe.Tire.InsertWord(keyword)
		ids = append(ids,id)
	}

	err := ldb.Put(k,utils.Encoder(ids))
	if err != nil {
		return
	}
}

//添加正排索引

// id(uint32) ：keywords map[string] []int

func (pe *PictureEngine) AddIndexOfIdKey(id uint32,keymap map[string] []int){

	 pe.Lock()
	 defer pe.Unlock()

	 ldb := pe.IdMapKey[pe.GetLeveldbId(id)]

	 k := utils.Uint32ToBytes(id)

	 flag, _ := ldb.Has(k)

	 if !flag {
		 err := ldb.Put(k,utils.Encoder(&keymap))
		 if err != nil {
			 return
		 }
	 }

}

// id (uint32) : Picture
func (pe *PictureEngine) AddDocument(id uint32,picture *Picture){
	 pe.Lock()
	 defer pe.Unlock()

	 ldb := pe.IdMapDocument[pe.GetLeveldbId(id)]

	k := utils.Uint32ToBytes(id)

	flag, _ := ldb.Has(k)

	if !flag {
		err := ldb.Put(k,utils.Encoder(picture))
		if err != nil {
			return
		}
	}
}

//BM25
//             IDF * TF * (k1 + 1)
//BM25 = sum ----------------------------
//            TF + k1 * (1 - b + b * D / L)
//
//                   总文档数目
//IDF = log2( ------------------------  + 1 )
//                出现该关键词的文档数目


func (pe *PictureEngine) GetRank(keywords []string, words map[string] []int,fre []uint32)float32{
	  IDF := pe.bm25.GetIDF(fre)
	  return pe.bm25.GetScore(keywords,words,IDF)
}









