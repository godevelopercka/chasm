package main

import "fmt"

func NewUser() {
	// 初始化结构体
	u := User{}
	fmt.Printf("%v \n", u)
	fmt.Printf("%+v \n", u) // 正常用这个仔细一点

	// up 是一个指针，指向 User 的结构体
	up := &User{}
	fmt.Printf("%+v \n", up)
	// 同 up
	up2 := new(User)
	fmt.Printf("%+v \n", up2)

	u4 := User{Name: "Tom", Age: 0}    // 用这种
	u5 := User{"Tom", "FirstName", 18} // 不要用这种，容易出错

	u4.Name = "Jerry"
	u5.Age = 18

	var up3 *User
	// nil 上访问字段、方法
	println(up3.FirstName) // 空指针：invalid memory address or nil pointer dereference
	println(up3)
}

type User struct {
	Name      string
	FirstName string
	Age       int
}

func (u *User) ChangeName(name string) {
	u.Name = name
}
func (u *User) ChangeAge(age int) {
	u.Age = age
}
