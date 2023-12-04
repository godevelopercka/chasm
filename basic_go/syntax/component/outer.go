package main

type Inner struct {
}

func (i Inner) DoSomething() {
	println("这是 Inner")
}

type Outer struct { // 正常用这个组合
	Inner
}

type OuterV1 struct { // 正常用这个组合
	Inner
}

func (o OuterV1) DoSomething() {
	println("这是 OuterV1")
}

type OOuter struct {
	Outer
}
type OuterPtr struct { // 一般不用这个
	*Inner
}
type OOOOuter struct {
	OOuter
}

func UseInner() {
	var o Outer
	o.DoSomething()       // 如果自己没有这个同名方法，组合了会先直接调用组合的方法
	o.Inner.DoSomething() // 调用组合的方法

	var op *OuterPtr
	op.DoSomething()

	//o2 := Outer{
	//	Inner: Inner{}, // 初始化组合
	//}
	//op2 := OuterPtr{
	//	Inner: &Inner{}, // 初始化组合
	//}
	//println(o2, op2) // 报错，不能 println 结构体
}

func (o Outer) Name() string {
	return "Outer"
}

func (i Inner) SayHello() { // 此时没有同名方法，所以这里的 i 调用 Inner 的 i
	println("hello," + i.Name())
}

func (o Outer) SayHello() { // 此时没有同名方法，所以这里的 i 调用 Inner 的 i
	println("hello," + o.Name())
}

func (i Inner) Name() string {
	return "Inner"
}

func UseOuter() {
	var o Outer
	o.SayHello()
}

func main() {
	o1 := OuterV1{}
	o1.DoSomething()       // 调用自己的实现
	o1.Inner.DoSomething() // 调用组合的实现
	UseOuter()
}
