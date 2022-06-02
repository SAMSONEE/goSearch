package trie

import (
	"SearchEngine/utils"
	"fmt"
)


func Serialize(node *trieNode, level rune,filename string)  {

	data1 := make([]rune,0)

	tmpnode1 := make([]*trieNode,0)

	data1 = append(data1,level)
	//第一个文件： 层数|节点|
	for word, nxt := range node.Children {
		data1 = append(data1,word)
		tmpnode1 = append(tmpnode1,nxt)
	}
	fmt.Println(data1)
	utils.Write(&data1, filename)

	//叶子节点 id 孩子数 节点 是否是词|
	for i := rune(1); i <= level; i++ {
		data :=  make([]rune,0)
		tmpnode2 := tmpnode1
		tmpnode1 = []*trieNode {}
		for id, children := range tmpnode2 {
			size := len(children.Children)
			data = append(data,rune(id))
			data = append(data,rune(size))
			for word, nxt := range children.Children {
				data= append(data,word)
				if children.Children[word].isWorld {
					data = append(data,rune(1))
				}else{
					data = append(data,rune(0))
				}
				tmpnode1 = append(tmpnode1,nxt)
			}
		}
		if len(tmpnode1) > 0 {
			fmt.Println(data)
			utils.Write(&data,fmt.Sprintf("%s%d",filename,i))
		}
	}

}


func Write(tree *TrieTree, filename string) {
	Serialize(tree.Root,tree.Level,filename)
}

func Read(filename string) *TrieTree {

	 tree := NewTrie()
	 data := make([]rune,0)
	 tmpnode := make([]*trieNode,0)
	 //解析第一个文件
	 utils.Read(&data, filename)
	 for id, word := range data {
		 if id == 0 {
			 tree.Level = word
		 }else {
			 tree.Root.Children[word] = NewNode()
			 tmpnode = append(tmpnode,tree.Root.Children[word])
		 }
	 }

	 for index := rune(1); index < tree.Level; index++ {
		 data1 := make([]rune,0)
		 tmpnode1 := tmpnode
		 tmpnode = []*trieNode {}
		 utils.Read(&data1, fmt.Sprintf("%s%d",filename,index))
		 dataSize := len(data1)
		 i := 0
		 for {
			 id := data1[i]
			 i++
			 childrensize := data1[i]
			 i++
			 //判断是否是最后一个元素
			 if  i >= dataSize {
				 break
			 }

			 for j := rune(0); j < childrensize; j++ {
				 word := data1[i]
				 i++
				 flag := data1[i]
				 i++
				 tmpnode1[id].Children[word] = NewNode()
				 if flag==rune(1) {
					 tmpnode1[id].Children[word].isWorld = true
				 }
				 tmpnode = append(tmpnode,tmpnode1[id].Children[word])
			 }
			 if  i >= dataSize {
				 break
			 }
		 }
	 }

	 return tree
}