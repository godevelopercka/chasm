package main

type LinkedList struct {
	head *node // 声明结构体字段的时候，需要用指针，这里的node是一个结构体，*node是获取node的指针，与node不一样
	tail *node
	Name string // 声明普通字段

	// 这个就是包外可访问
	Len int
}

// 私有结构体
type node struct {
}

func (l LinkedList) Add(idx int, val any) error {
	//TODO implement me
	panic("implement me")
}

func (l LinkedList) Append(val any) {
	//TODO implement me
	panic("implement me")
}

func (l LinkedList) Delete(index int) (any, error) {
	//TODO implement me
	panic("implement me")
}

func (l LinkedList) toSlice() ([]any, error) {
	//TODO implement me
	panic("implement me")
}

// 结构体接收器
//func (l LinkedList) Add(idx int, val any) error {
//	//TODO implement me
//	panic("implement me")
//}
// 方法接收器，receiver
//func (l *LinkedList) AddV1(idx int, val any) {
//
//}
