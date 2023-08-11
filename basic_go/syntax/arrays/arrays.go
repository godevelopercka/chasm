package main

//func array() {
//	// 直接初始化一个三个元素的数组，大括号里面只能少，不能多
//	arr1 := [3]int{1, 2, 3}
//	fmt.Printf("元素：%v, 数组目前所占长度：%d, 数组最大容量(即最长长度)：%d", arr1, len(arr1), cap(arr1))
//
//	// 少了的部分就是默认零值，等价于 3, 2, 0
//	arr2 := [3]int{3, 2}
//	fmt.Printf("arr2: %v, len: %d, cap: %d", arr2, len(arr2), cap(arr2))
//
//	// 虽然没有显示初始化，但实际上内存已经分配好，等价于 0, 0, 0
//	var arr3 [3]int
//	fmt.Printf("arr2: %v, len: %d, cap: %d", arr3, len(arr3), cap(arr3))
//
//	// 数组不支持 append 操作
//	// arr3 = append(arr3, 1)
//
//	// 按下标索引，如果编译器能判断出下标越界，那么会编译错误
//	// 如果不能，那么运行时候会报错，出现 panic
//	fmt.Printf("arr1[1]:%d", arr1[1])
//
//}
//
//func Slice() {
//	s1 := []int{1, 12, 3, 5} // 直接初始化了 4 个元素的切片
//	fmt.Printf("s1: %v, len: %d, cap: %d \n", s1, len(s1), cap(s1))
//
//	s2 := make([]int, 3, 4) // 直接初始化了三个元素，容量为 4 的切片
//	fmt.Printf("s2: %v, len: %d, cap: %d \n", s2, len(s2), cap(s2))
//
//	s2 = append(s2, 3)                                              // 再追加一个元素，没有扩容
//	s2 = append(s2, 8)                                              // 再追加一个元素，扩容了
//	s3 := make([]int, 4)                                            //make 只传入一个参数，表示创建一个 4 个元素的切片
//	fmt.Printf("s3: %v, len: %d, cap: %d \n", s3, len(s3), cap(s3)) // 输出 s3: [0 0 0 0], len: 4, cap: 4
//
//	// 按照下标索引
//	fmt.Printf("s3[2]:%d", s3[2])
//	// 超出下标返回, panic
//	fmt.Printf("s3[2]:%d", s3[99])
//}
//func Map1() {
//	m1 := map[string]string{
//		"key1": "value1",
//	}
//
//	val, ok := m1["key1"]
//	println(val, ok)
//	if ok {
//		println("第一步:", val)
//	}
//	val = m1["key2"]
//	println("第二步:", val)
//}
//
//func Map2() {
//	m2 := make(map[string]string, 4)
//	m2["key2"] = "value2"
//
//	println(len(m2))
//	for k, v := range m2 {
//		println(k, v)
//	}
//	for k := range m2 {
//		println(k)
//	}
//
//	delete(m2, "key2")
//}

// 第二讲课后作业（2）:获得 map 的所有 key、所有 value
func MapReturn(s map[string]string) (k string, v string) {
	for k, v := range s {
		println(k, v)
	}
	return k, v
}
