package startup

import (
	"webook_go/webook/pkg/logger"
)

func InitLog() logger.LoggerV1 {
	return &logger.NopLogger{}
}
