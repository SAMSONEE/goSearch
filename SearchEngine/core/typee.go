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

type Searchrequest struct {
	QueryText		 string
	KeyWords	 	string
	HateWords	 	string
}

type Serachresult struct {
	//Time 			float32
	Total 			int
	PageCount		int				//总页数
	Page 			int				//页码
	Limit 			int				//页大小
	Pictures		[]Picture	//存储信息
}

type Score struct {
	Id 	uint32
	Score float32
}

