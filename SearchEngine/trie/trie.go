//字典树
//存在的问题
//后面优化

package trie

//字典树节点
type trieNode struct {
	 isWorld bool
	 Children map[rune] *trieNode
}

//字典树
type TrieTree struct {
	 Level rune
	 Root *trieNode
}

//创建新节点
func NewNode() *trieNode {
	return & trieNode{isWorld: false,Children: make(map[rune]*trieNode)}
}

//创建字典树
func NewTrie() *TrieTree {
	 return &TrieTree{
		 Level:0,
		 Root :&trieNode{isWorld: false,Children: make(map[rune]*trieNode)},
	 }
}


//可以优化
func (t *TrieTree) InsertWord(word string)  {
	 if len(word) == 0 {
		 return
	 }

	level := t.Level
	i := rune(0)

	root  := t.Root

	for _, w := range word {
		 //是否存在节点，存在新建一个
		 if _, ok := root.Children[w]; !ok {
			 	root.Children[w]=NewNode()
		 }
		 root = root.Children[w]
		 i++
	 }
	 if level < i {
		 t.Level =i
	 }
	 if !root.isWorld {
		 root.isWorld = true
	 }
}

//是否包含单词
func (t *TrieTree) Contains(word string) bool {
	 if len(word) == 0 {
		 return false
	 }
	 root := t.Root
	for _, w :=range word {
		 if node, ok := root.Children[w]; ok{
			 root = node
		 }else{
			 return false
		 }
	 }
	 return root.isWorld
}


//是否是前缀
//可用来相关搜索

func (t *TrieTree) IsPrefix(word string) bool{
	if len(word) == 0 {
		return false
	}

	root := t.Root
	for _, w := range word {
		if node, ok := root.Children[w]; ok{
			root = node
		} else {
			return false
		}
	}
	return true
}






