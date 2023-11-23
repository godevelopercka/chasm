package main

import "fmt"

func Sum[T Number](vals ...T) T {
	var res T
	for _, val := range vals {
		res = res + val
	}
	return res
}

func Mul[T Number](res T, vals ...T) T {
	for _, val := range vals {
		res = res * val
	}
	return res
}

func Div[T Number](res T, vals ...T) T {
	for _, val := range vals {
		res = res / val
	}
	return res
}

type Number interface {
	~int | int32 | int64 | float32 | float64
}

func main() {
	//fmt.Println(Sum[float32](1.12, 21.231123, 1.123))
	//fmt.Println(Mul[float32](6.0, 2.0))
	fmt.Println("订单转化率-SUM(广告订单量) / SUM(点击量：)", Div[float32](6, 69))
	fmt.Println("ACOS-SUM(广告花费) / SUM(广告销售额：)", Div[float32](47.92, 156.94))
}
