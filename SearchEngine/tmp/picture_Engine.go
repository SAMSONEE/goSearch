package main

import (
	"github.com/huichen/sego"
	"log"
)

type pictureEngine struct {
	//搜索信息
	pictures []Picture
	//停字符
	stoptokens StopTokens
	//分词器
	segmenter sego.Segmenter
}

func (pe *pictureEngine) Init() {
	//加载词典
	pe.segmenter.LoadDictionary("data/dictionary.txt")

	//通用词路径
	stopname := "./data/stop_tokens.txt"
	//初始化停用词
	pe.stoptokens.ReadStop(stopname)
}

//对csv进行预处理
func (pe *pictureEngine) SegmentCsv(filePath string) {

	//读取CSV文件
	message, err := Readcsv(filePath)

	if err != nil {
		log.Fatalln(err)
	}

	//如何将string 转化为byte
	for _, t := range message {

		tokensmap := make(map[string][]int)

		text := []byte(t[1])
		segments := pe.segmenter.Segment(text)
		// words := sego.SegmentsToSlice(segments, true)

		for _, segment := range segments {
			token := segment.Token().Text()

			if !pe.stoptokens.Stop_Tokens[token] {
				tokensmap[token] = append(tokensmap[token], segment.Start())
			}
		}

		contextlen := len(segments)

		pe.pictures = append(pe.pictures, Picture{t[0], tokensmap, contextlen,t[1]})

	}
}
