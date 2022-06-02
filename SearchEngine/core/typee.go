package core

//停止字符
//在分词操作中去除这些单词
type StopTokens struct {
	Stop_Tokens map[string]bool
}

//存储数据集信息
type Picture struct {
	Id 					uint32
	//照片url
	Picture_url         string
	//文本内容
	Picture_context     string
}