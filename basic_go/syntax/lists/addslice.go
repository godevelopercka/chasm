package main

import "errors"

func AddSlice[T any](slice []T, idx int, val T) ([]T, error) {
	// 如果我这边 idx 是负数，或者超过了 slice 的界限
	if idx < 0 || idx > len(slice) {
		return nil, errors.New("下标出错")
	}

	res := make([]T, 0, len(slice)+1)
	for i := 0; i < idx; i++ {
		res = append(res, slice[i])
	}

	res[idx] = val
	for i := idx; i < len(slice); i++ {
		res = append(res, slice[i])
	}

	return res, nil
}

func main() {
	println(AddSlice[int]([]int{12, 2, 4}, 6, 2))
}
