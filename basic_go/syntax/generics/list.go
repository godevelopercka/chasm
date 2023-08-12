package main

// T 类型参数，名字叫做 T，约束是 any，等于没有约束
type List[T any] interface {
	Add(idx int, t T)
	Append(t T)
}

func main() {
	//UseList()                       // 运行会报空指针，因为没有初始化
	println(Sum[int](1, 2, 3))      // 只能计算同类型
	println(Sum[Interger](1, 3, 5)) // 衍生类型，~int
	println(Sum[float64](1.12, 21.231123, 1.123))

}

func UseList() {
	var l List[int]
	l.Append(18)

	var j List[any]
	j.Append("string")
	lk := LinkedList[int]{} // 初始化泛型结构体
	intVal := lk.head.val
	println(intVal)
}

type LinkedList[T any] struct {
	head *node[T]
	t    T // 声明一个 T 类型的字段
}

type node[T any] struct {
	val T
}
