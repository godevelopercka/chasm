package wire

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"webook_go/wire/repository"
	"webook_go/wire/repository/dao"
)

func main() {
	db, err := gorm.Open(mysql.Open("dsn"))
	if err != nil {
		panic(err)
	}
	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	fmt.Println(repo)
}
