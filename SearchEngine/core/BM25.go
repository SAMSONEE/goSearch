package core


// BM25算法用来衡量一组关键词和某文档的相关程度

/*

                IDF * TF * (k1 + 1)
BM25 = sum ----------------------------
           TF + k1 * (1 - b + b * D / L)

 */


/*
*   TF： 某关键词的词频
*   D : 该文档的词数
*   L : 所有文档词数的平均数
*   k1 = 2.0 b = 0.75 默认值
*
*  IDF（inverse document frequency）衡量关键词是否常见
 */

/*
                   总文档数目
IDF = log2( ------------------------  + 1 )
             出现该关键词的文档数目
 */

type BM25Parameter struct{
	 K1 float32
	 B  float32
}

