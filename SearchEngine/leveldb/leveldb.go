/*
*   Todo：
*	     1.安全性问题
*        2. 数据压缩
*
 */


package leveldb

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type Leveldb struct {
	 db *leveldb.DB
	 path string
}

func CreateLeveldb(path string) (*Leveldb){
	db, err := leveldb.OpenFile(path,nil)

	//fmt.Println(err)

	//错误没写

	if err != nil {
		return nil
	}

	result := &Leveldb{
		db: db,
		path: path,
	}

	return result
}

//string to []byte
func (db *Leveldb) Get(key []byte)([]byte, error){
	return db.db.Get(key,nil)
}

func(db *Leveldb) Put(key []byte, val []byte) error{
	 return db.db.Put(key,val,nil)
}

func(db *Leveldb) Has(key []byte) (bool,error){
	return db.db.Has(key,nil)
}

func(db *Leveldb) Delete(key []byte) error {
	return db.db.Delete(key,nil)
}

func(db *Leveldb) Close()error {
	return db.db.Close()
}
