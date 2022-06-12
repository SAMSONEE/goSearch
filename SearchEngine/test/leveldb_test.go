package test

import (
	"SearchEngine/leveldb"
	"SearchEngine/utils"
	"fmt"
	"math/rand"
	"testing"
	"time"
)


//func DB(){
//
//	filepath := "/home/wqk/GolandProjects/SearchEngine/test/testLeveldb"
//
//	ldb := leveldb.CreateLeveldb(filepath)
//
//	for i := uint32(0); i < uint32(10); i++ {
//
//		k := utils.Uint32ToBytes(i)
//
//		err := ldb.Put(k,k)
//		if err != nil {
//			return
//		}
//		fmt.Println(k)
//	}
//
//	rand.Seed(time.Now().UnixNano())
//
//	for i := 0; i < 10; i++ {
//		randid := uint32(rand.Intn(10))
//
//		k := utils.Uint32ToBytes(randid)
//
//		v, _ := ldb.Get(k)
//
//		fmt.Println(utils.BytesToUint32(v))
//	}
//}

func DB2() {
	filepath := "/home/wqk/GolandProjects/SearchEngine/test/testLeveldb"

	ldb2, err := leveldb.Open(filepath)

	if err != nil {
		fmt.Println(err)
		return
	}

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 10; i++ {
		randid := uint32(rand.Intn(10))

		k := utils.Uint32ToBytes(randid)

		v, _ := ldb2.Get(k)

		fmt.Println(utils.BytesToUint32(v))
	}

}


//func TestDB(t *testing.T) {
//	 DB()
//}

func TestDB2(t *testing.T) {
	DB2()
}