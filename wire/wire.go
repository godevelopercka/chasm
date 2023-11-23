//go:build wireinject

// 让 wire 来注入这里的代码

package wire

import (
	"github.com/google/wire"
	"practice/wire/repository"
	"practice/wire/repository/dao"
)

func InitRepository() *repository.UserRepository {
	// 我只在这里声明我要用的各种东西，但是具体怎么构造，怎么编排排序
	// 这个方法里面传入各个组件的初始化方法
	wire.Build(repository.NewUserRepository, dao.NewUserDAO, InitDB)
	return new(repository.UserRepository)
}
