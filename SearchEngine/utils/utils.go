/*
*   "encoding/gob" 原理
*
*
*
*/

package utils

import (
	"bytes"
	"compress/flate"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"log"
)

//编码
func Encoder(data interface{}) []byte {
	 if data == nil {
		 return nil
	 }

	 buffer := new(bytes.Buffer)
	 //编码后数据写入buffer
	 encoder := gob.NewEncoder(buffer)

	 err := encoder.Encode(data)

	 if err != nil {
		 log.Fatal("Encoder:", err)
	 }
	 return buffer.Bytes()
}

//解码
func Decoder(data []byte, v interface{}){
	if data == nil {
		return
	}
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(v)
	if err != nil {
		log.Fatal("Decoder:", err)
	}
}

// StringToInt 字符串转整数
func StringToInt(value string) uint32 {
	//crc32 算法
	return crc32.ChecksumIEEE([]byte(value))
}


func Uint32ToBytes(i uint32) []byte {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, i)
	return buf
}


func BytesToUint32(buf []byte) uint32 {
	return binary.BigEndian.Uint32(buf)
}

//写盘操作
func Write(data interface{}, filename string) {
	 buffer := new(bytes.Buffer)
	 encoder := gob.NewEncoder(buffer)

	 err := encoder.Encode(data)
	 if err != nil {
		 log.Println(err)
		 return
	 }

	 log.Println("写入文件：",filename)
	 compressData := compression(buffer.Bytes())

	 err = ioutil.WriteFile(filename,compressData,0600)
	 if err != nil {
		 log.Println(err)
		 return
	 }

}

// 文件压缩 flate 可用其他的
func compression(data []byte) []byte {

	buf := new(bytes.Buffer)
	write, err := flate.NewWriter(buf,flate.DefaultCompression)
	defer write.Close()

	if err != nil {
		log.Println(err)
		return nil
	}

	write.Write(data)
	write.Flush()
	log.Println("原大小：", len(data), "压缩后大小：", buf.Len(), "压缩率：", fmt.Sprintf("%.2f", float32(buf.Len())/float32(len(data))), "%")
	return buf.Bytes()
}

//Decompression
func Decompression(data []byte) []byte {
	buf := new(bytes.Buffer)
	read := flate.NewReader(bytes.NewReader(data))
	defer read.Close()

	buf.ReadFrom(read)

	return buf.Bytes()
}

func Read(data interface{},filename string){
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return
	}

	decoData := Decompression(raw)

	buffer := bytes.NewBuffer(decoData)
	dec := gob.NewDecoder(buffer)
	err = dec.Decode(data)
	if err != nil {
		log.Println(err)
		return
	}
}