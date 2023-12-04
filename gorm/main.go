package main

import "fmt"

//
//import (
//	"gorm.io/driver/mysql"
//	"gorm.io/gorm"
//)
//
//type Product struct {
//	gorm.Model
//	Code  string `gorm:"unique"`
//	Price uint
//}
//
//func main() {
//	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13317)/webookv1"))
//	if err != nil {
//		panic("failed to connect database")
//	}
//
//	db = db.Debug()
//	// 迁移 schema
//	// 建表
//	db.AutoMigrate(&Product{})
//
//	// Create
//	db.Create(&Product{Code: "D42", Price: 100})
//
//	// Read
//	var product Product
//	db.First(&product, 1)                 // 根据整型主键查找
//	db.First(&product, "code = ?", "D42") // 查找 code 字段值为 D42 的记录
//
//	// Update - 将 product 的 price 更新为 200
//	db.Model(&product).Update("Price", 200)
//	// Update - 更新多个字段
//	db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
//	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})
//
//	// Delete - 删除 product
//	db.Delete(&product, 1)
//}

type MyStruct struct {
	name string
}

func Fun1() *MyStruct {
	a := &MyStruct{
		name: "SuXi",
	}
	defer func() {
		a.name = "linxiaoyang"
	}()
	return a
}

func main() {
	fmt.Println(Fun1())
	fmt.Println(MyStruct{})
}
