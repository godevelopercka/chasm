package main

type List interface {
	Add(idx int, val any) error
	Append(val any)
	Delete(index int) (any, error)
	toSlice() ([]any, error) // 包外无法调用
}
