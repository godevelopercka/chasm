package main

import (
	"io"
)

// 泛型的 sum 求和方法
func Sum[T Number](vals ...T) T {
	var res T
	for _, val := range vals {
		res = res + val
	}
	return res
}

type Number interface {
	~int | int32 | int64 | float32 | float64
}

type Interger int

// func SumV1[T Number](vals ...Number) Number { // Number 约束不能用于参数和返回值，需改成 T
//
//		var t T
//		return t
//	}
func SumV2[T Number](vals ...T) T { // Number 约束不能用于参数和返回值，需改成 T
	var t T
	return t
}
func ReleaseResource[R io.Closer](r R) { // 约束要么是any，要么是接口
	r.Close()
}
