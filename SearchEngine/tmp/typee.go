package main

//停止字符
//在分词操作中去除这些单词
type StopTokens struct {
	Stop_Tokens map[string]bool
}

//存储数据集信息
type Picture struct {
	Picture_url         string
	Picture_tokens_map  map[string][]int
	Picture_context_len int
	Picture_context string
}
