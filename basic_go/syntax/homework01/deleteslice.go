package main

import (
	"errors"
	"fmt"
)

func DeleteSlice[T any](slice []T, idx int) ([]T, error) {
	if idx < 0 || idx >= len(slice) {
		return nil, errors.New("下标错误")
	}
	// 定义一个切片
	res := make([]T, 0, len(slice))
	for i := 0; i < idx; i++ {
		res = append(res, slice[i])
	}
	for i := idx + 1; i < len(slice); i++ {
		res = append(res, slice[i])
	}
	// 缩容
	ret := make([]T, 0, len(res))
	for i := 0; i < len(res); i++ {
		ret = append(ret, res[i])
	}
	return ret, nil
}
func main() {
	res, _ := DeleteSlice[int]([]int{1, 23, 45, 62, 3}, 2)
	fmt.Println(res)
	//fmt.Println(DeleteSlice[int]([]int{1, 23, 45, 62, 3}, 5))
	fmt.Println(len(res), cap(res))
}
