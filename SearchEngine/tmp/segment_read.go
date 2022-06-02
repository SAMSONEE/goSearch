package main

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
)


//初始化 停止符
func (st *StopTokens) ReadStop(filepath string) {

	st.Stop_Tokens = make(map[string]bool)

	if filepath == "" {
		return
	}
	file, err := os.Open(filepath)

	if err != nil {
		log.Println(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		text := scanner.Text()

		if text != "" {
			st.Stop_Tokens[text] = true
		}
	}
}

// 功能：读取Csv文件
// 输入： 文件路径名
// 输出： 二维sring [url context]
func Readcsv(filepath string) ([][]string, error) {

	opencast, err := os.Open(filepath)

	if err != nil {
		log.Println("csv文件打开失败!")
	}

	defer opencast.Close()

	Readcsv := csv.NewReader(opencast)

	//去除第一行 标签栏 url contex
	Readcsv.Read()

	//读取一行数据
	// read, err := Readcsv.Read()
	read, err := Readcsv.ReadAll()

	return read, err
}
