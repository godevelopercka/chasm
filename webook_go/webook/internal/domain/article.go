package domain

// Article 可以同时表达线上库和制作库的概念吗?
// 可以同时表达，作者眼中的 Article 和读者眼中的 Article 吗?
type Article struct {
	Id      int64
	Title   string
	Content string
	// Author 要从用户来
	Author Author
}

// Author 在帖子这个领域内，是一个值对象
type Author struct {
	Id   int64
	Name string
}
