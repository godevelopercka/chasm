package main

func Func0(name string) string {
	return "hello, world"
}

// Func1 多个返回值
func Func1(name, hobby, gender string, age int) (string, error) {
	return "", nil
}

// Func2 带名字的返回值
func Func2(name string, age int) (str string, err error) {
	str = "hello"
	return
}

// Func3 带名字的返回值
func Func3(a int, b int) (name string, age int) {
	res := "小明" // 新的局部变量才能使用
	age = 18
	// 虽然带名字，但我们并没有用
	return res, age
}

func Invoke() {
	str, err := Func2("小明", 18)
	println(str, err)
	// 忽略返回值
	_, _ = Func2("肥猫", 6)
	// 忽略部分返回值
	// str是已经声明好了
	str, _ = Func2("肥仔", 20)
	// str1 是新变量，需要使用 :=
	str1, _ := Func2("肥佬", 21)
	println(str1)
	// str2是新变量，需要使用 :=
	str2, err := Func2("肥牛", 10)
	println(str2)

}

// Recursive 递归
// 这个方法运行的时候会出现错误
//func Recursive() {
//	Recursive()
//}

func Func4() {
	myFunc3 := Func3
	_, _ = myFunc3(1, 2)
}

func Func5() {
	fn := func(name string) string {
		return "hello, " + name
	}
	fn("大明")
}

func Func7() func(name string) string { // func(name string)相当于返回值名称
	return func(name string) string {
		return "hello," + name
	}
}

func Func8(str string) {
	hello := func(name string) string {
		return "hello, world" + name
	}
	println(hello(str))
}

func Closure1(name string) func() string {

	// 返回的这个函数，就是一个闭包。
	// 它引用到了 Closure 这个方法的入参
	return func() string {
		return "hello," + name
	}
}

func Func9(name string, aliases ...string) {
	if len(aliases) > 0 {
		println(aliases[0])
	}
}

func UseFunc9() {
	Func9("xiaoming")
	Func9("xiaoming", "daming")
	Func9("xiaoming", "daming", "zhongming")
}
