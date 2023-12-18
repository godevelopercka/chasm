package main

import (
	"webook_go/webook/internal/service/sms"
	"webook_go/webook/internal/service/sms/memory"
)

func InitSMSService() sms.Service {
	// 换内存还是换第三方
	return memory.NewService()
}
