package startup

import "webook_go/webook/internal/web"

func InitWechatHandlerConfig() web.WechatHandlerConfig {
	return web.WechatHandlerConfig{
		Secret: false,
	}
}
