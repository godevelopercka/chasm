package main

func Defer() {
	defer func() {
		println("第一个defer")
	}()

	defer func() {
		println("第二个defer")
	}()
}

func Defer1() {
	i := 0
	defer func() {
		println(i)
	}()
	i = 1
}

func Defer2() {
	i := 0
	defer func(val int) {
		println(val)
	}(i)
	i = 1
}

func DeferReturn() int {
	a := 0
	defer func() {
		a = 1
	}()
	return a
}

func DeferReturnV1() (a int) {
	a = 0
	defer func() {
		a = 1
	}()
	return a
}

func DeferReturnV2() *MyStruct {
	a := &MyStruct{
		name: "Jerry",
	}
	defer func() {
		a.name = "Tom"
	}()
	return a
}

type MyStruct struct {
	name string
}

func IfOnly(biu int, age int) string {
	if age >= 18 {
		return "成年了"
	} else if age > 12 {
		return "骚年"
	} else {
		return "他还是个孩子"
	}
}

func IfNewVariable(start int, end int) string {
	if distance := start - end; distance > 100 {
		println(distance)
		return "距离太远了"
	} else {
		println(distance)
		return "距离比较近"
	}
}

func Loop1() {
	for i := 0; i < 8; i++ {
		println(i)
	}

	// 这样也可以
	for i := 0; i < 9; {
		println(i)
		i++
	}
}

func Loop2() {
	i := 0
	for i < 9 {
		println(i)
		i++
	}
}
func Loop3() {
	for {
		println("hello")
	}
}

func Loop4() {
	arr := [3]int{1, 2, 3}
	for k, v := range arr {
		println(k, v)
	}
	for i := range arr {
		println(i, arr[i])
	}
}

func Loop5() {
	m := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	for k, v := range m {
		print(k, v)
	}
	for i := range m {
		println(i, m[i])
	}
}

func Loop6() {
	i := 0
	for {
		if i > 2 {
			println("已中止")
			break
		}
		i++
	}
}
func Loop7() {

	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			println("已跳出当前循环，开始下次循环")
			continue
		}

	}
}

func Switch(status int) string {
	switch status {
	case 0:
		return "初始化"
	case 1:
		return "运行中"
	default:
		return "default 分支可有可无"
	}
}

func Switch1(age int) {
	switch { // switch 语句无value时，case 后面跟 bool 表达式，且每一个条件要互斥
	case age > 18:
		println("成年人")
	case age <= 18:
		println("小孩子")
	}
}
