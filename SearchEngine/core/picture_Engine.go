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
	"math"
	"sync"
)



type PictureEngine struct {

	//被索引的文档数
	numDocumentIsIndex uint32
	//被分词的文档数
	numSegment uint32
	////被存储的文档数
	//numDocumentStored uint32

	////被存储的文档的平均词数
	//numAvergeword uint8

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
	fileTrie := "/home/wqk/GolandProjects/SearchEngine/test/data/storetrie/st"
	fileBM25 := "/home/wqk/GolandProjects/SearchEngine/test/data/BM25/BM"
	fileElse := "/home/wqk/GolandProjects/SearchEngine/test/data/else/else"



	//初始化leveldb
	for i := 0 ; i < 10; i++ {

		Kmi , _ := leveldb.Open(fmt.Sprintf("%s%d",Keypath,i))
		pe.KeyMapId = append(pe.KeyMapId,Kmi)

		Imk , _ := leveldb.Open(fmt.Sprintf("%s%d",Idpath,i))
		pe.IdMapKey = append(pe.IdMapKey,Imk)

		Imd , _ := leveldb.Open(fmt.Sprintf("%s%d",Docpath,i))
		pe.IdMapDocument = append(pe.IdMapDocument,Imd)

		//pe.KeyMapId = append(pe.KeyMapId,leveldb.CreateLeveldb(fmt.Sprintf("%s%d",Keypath,i)))
		//pe.IdMapKey = append(pe.IdMapKey,leveldb.CreateLeveldb(fmt.Sprintf("%s%d",Idpath,i)))
		//pe.IdMapDocument = append(pe.IdMapDocument,leveldb.CreateLeveldb(fmt.Sprintf("%s%d",Docpath,i)))
	}

	//初始化字典树
	pe.Tire = Read(fileTrie)

	p, total, values := pe.bm25.ReadBM25(fileBM25)

	ma := movingaverage.New(math.MaxUint16)

	pa := &rank.BM25Parameter{
		K1: p[0],
		B: p[1],
	}

	ma.Addnums(values)

	pe.bm25 = &rank.BM25{
		Parameters: pa,
		Total: total,
		Ma: ma,
	}

	es := make([]uint32,0)

	utils.Read(&es,fileElse)

	pe.numDocumentIsIndex = es[0]
	pe.numSegment = es[1]

	fmt.Println("初始化完成")
}

func (pe* PictureEngine) Debug()  {
	fmt.Println("**************************")
	fmt.Println(pe.bm25.Parameters)
	fmt.Println("Total:",pe.bm25.Total)
	fmt.Println("values:")
	fmt.Println(pe.bm25.Ma.GetValues())
	fmt.Println(pe.numDocumentIsIndex," ",pe.numSegment)
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


//多通道搜索
func (pe *PictureEngine) Search(request Searchrequest) []Serachresult {

	querytext := request.QueryText
	hatetext := request.HateWords
	keytext := request.KeyWords

	//等待初始化完成
	//pe.Wait()

	//分词
	text := []byte(querytext+keytext)
	segments := pe.segmenter.Segment(text)
	words := make(map[string] bool)
	for _, segment := range segments {
		token := segment.Token().Text()

		//不是分词符号则加入
		if !pe.stoptokens.Stop_Tokens[token] {
			if !words[token] {
				words[token] = true
			}
		}
	}

	//屏蔽词
	text2 := []byte(hatetext)
	segments2 := pe.segmenter.Segment(text2)
	hatemap := make(map[string] bool)

	for _, segment := range segments2 {
		token := segment.Token().Text()

		//不是分词符号则加入
		if !pe.stoptokens.Stop_Tokens[token] {
			if !hatemap[token] {
				hatemap[token] = true
			}
		}
	}


	keywords := make([]string,0)

	for item, _ := range words {
		if !hatemap[item] {
			keywords = append(keywords, item)
		}
	}

	fmt.Println("keywords:",keywords)
	//fmt.Println("len:",len(keywords))

	fre := make([]uint32,0)
	//fmt.Println("fre:",fre)
	ids := make([]uint32,0)

	for _, item := range keywords {
		keyint := utils.StringToInt(item)
		k := utils.Uint32ToBytes(keyint)
		ldb := pe.KeyMapId[pe.GetLeveldbId(keyint)]
		if pe.Tire.Contains(item){
			flag, _ := ldb.Has(k)
			tmpid := make([]uint32,0)
			if flag {
				buf, _ := ldb.Get(k)
				utils.Decoder(buf,&tmpid)
			}
			//fmt.Println("tmpid:",tmpid)
			fre = append(fre,uint32(len(tmpid)))
			ids = append(ids,tmpid...)
		}
	}

	//fmt.Println("Ids:",len(ids))
	//fmt.Println("Frequence:",fre)
	//fmt.Println("len(words)",len(keywords))

	//多通道搜索

	fastsort := new(FastSort)
	//条件变量
	var wg sync.WaitGroup
	wg.Add(len(keywords)+1)

	//分组
	size := len(ids)/len(keywords)


	//fmt.Println("size:",size)
	//fmt.Println("hatewords:",hatewords)

	for i,j := 0,len(keywords); i <= j; i++ {
		go func(index int,s int,l int) {

			//fmt.Println("go func(index int,size int):",index)
			scs := make([]*Score,0)
			if s > 0 {
				if index != l {
					scs = pe.SimpleSearch(ids[index*s:(index+1)*s],keywords,hatemap,fre)
				}else {
					scs = pe.SimpleSearch(ids[index*s:],keywords,hatemap,fre)
				}
			}
			//fmt.Println("scs:",len(scs))
			wg.Done()
			//fmt.Println(scs)
			fastsort.Add(scs)
		}(i,size,j)
	}

	wg.Wait()
	sorted := fastsort.GetSort()


	total := len(sorted)
	//fmt.Println("total:",total)
	//默认每页20条记录

	size2 := total/20

	//fmt.Println("size2:",size2)

	result := make([]Serachresult,size2+1)

	wg.Add(size2+1)

	for i := 0; i <= size2 ; i++ {
		go func(index int, l int) {
			//fmt.Println("i:",index)
			pictures := make([]Picture,0)
			ids := make([]Score,0)
			if index != l {
				ids = sorted[index*20 : (index+1)*20]
			}else {
				ids = sorted[index*20:]
			}

			for _, item := range ids {
				tmpdb := pe.IdMapDocument[pe.GetLeveldbId(item.Id)]

				tmpid := utils.Uint32ToBytes(item.Id)

				flag, _ := tmpdb.Has(tmpid)

				if flag {
					var picture Picture

					buf, _ := tmpdb.Get(tmpid)
					utils.Decoder(buf,&picture)

					pictures = append(pictures,picture)
				}
			}

			result[index] = Serachresult{
				Total: total,
				PageCount: size2+1,
				Page: index+1,
				Limit: 20,
				Pictures: pictures,
			}
			wg.Done()
		}(i,size2)
	}

	wg.Wait()
	return result
}


//搜索

func (pe *PictureEngine) SimpleSearch(ids []uint32, keywords []string, hatewords map[string]bool,frequence []uint32) []*Score {
	Scs := make([]*Score,0)

	var wg sync.WaitGroup
	wg.Add(len(ids))
	//fmt.Println("ids:",len(ids))
	for _, i := range ids {
		go func(id uint32) {

			//fmt.Println("go func(id uint32):",id)

			ldb2 := pe.IdMapKey[pe.GetLeveldbId(id)]

			k := utils.Uint32ToBytes(id)

			flag, _ := ldb2.Has(k)

			if flag {
				keymap := make(map[string][]int,0)
				buf, _:= ldb2.Get(k)
				utils.Decoder(buf,&keymap)

				flag2 := false

				for item, _:= range keymap {
					if hatewords[item]  {
						flag2 = true
						break
					}
				}

				if !flag2 {
					score := pe.GetRank(keywords,keymap,frequence)
					//fmt.Println(id ," ",score)
					Sc := &Score{
						Id: id,
						Score: score,
					}
					//fmt.Println("Id:",id," Score:",score)
					Scs = append(Scs,Sc)
				}
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
	return Scs
}


//定时把字典树和BM25更新到磁盘
func (pe *PictureEngine) Flush() {
	pe.Lock()
	defer pe.Unlock()

	fileTrie := "/home/wqk/GolandProjects/SearchEngine/test/data/storetrie/st"

	Write(pe.Tire,fileTrie)

	fileBM25 := "/home/wqk/GolandProjects/SearchEngine/test/data/BM25/BM"

	pe.bm25.WriteBM25(fileBM25)

	//
	//
	//
	fileElse := "/home/wqk/GolandProjects/SearchEngine/test/data/else/else"

	tmps := []uint32{pe.numDocumentIsIndex,pe.numSegment}

	utils.Write(&tmps,fileElse)

}






