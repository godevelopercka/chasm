package main

import (
	"fmt"
	"time"
)

//type node struct {
//	prev       *node
//	next       *node
//	CreateTime time.Time // 引用其他包
//	// 自引用不用指针会编译错误
//	//next node
//}

func NewUser() {
	// 初始化结构体
	u := User{}
	fmt.Printf("%v \n", u)  // %v: 只能打印值 {  0}
	fmt.Printf("%+v \n", u) // %+v: 打印字段和值 {Name: FirstName: Age:0} 正常用这个仔细一点

	// up 是一个指针，指向 User 的结构体
	up := &User{}
	fmt.Printf("%+v \n", up) // &{Name: FirstName: Age:0}
	// 同 up
	up2 := new(User)
	println(up2.FirstName)
	fmt.Printf("%+v \n", up2) // &{Name: FirstName: Age:0}

	u4 := User{Name: "Tom", Age: 0}    // 用这种
	u5 := User{"Tom", "FirstName", 18} // 不要用这种，容易出错

	u4.Name = "Jerry"
	u5.Age = 18

	//var up3 *User
	// nil 上访问字段、方法
	//println(up3.FirstName) // 空指针：invalid memory address or nil pointer dereference
	//println(up3)

}

type User struct {
	Name      string
	FirstName string
	Age       int
}

func (u User) ChangeName(name string) {
	u.Name = name
}
func (u *User) ChangeAge(age int) {
	u.Age = age
}
func ChangeUser() {
	up4 := &User{Name: "xiaoming", FirstName: "xiaohei", Age: 18}
	up4.ChangeName("daming") // 结构接收器相当于复制体，无法修改 User 本体
	up4.ChangeAge(19)        // 指针接收器可以修改 User 本体
	//println(up4) // 结构体不能直接用 println 报错：illegal types for operand: print
	fmt.Printf("%+v", up4) // 输出 {Name:xiaoming FirstName:xiaohei Age:19}
}

type Interger int // Interger 是 int 的衍生类型

func UseInt() {
	i1 := 10
	i2 := Interger(i1)
	var i3 Interger = 11
	println(i2, i3) // 输出 10 11
}

type fish struct {
	Name string
}

func (f fish) Swim() {
	fmt.Printf("fish 在游")
}

type Fakefish fish

type Yu = fish

func UseFish() {
	f1 := fish{}
	f2 := Fakefish(f1)
	//f2.Swim()
	println(f2)
	y := Yu{}
	y.Swim()
}

type MyTime time.Time

type MyTim struct {
}

func (m MyTim) MyFunc() {

}
