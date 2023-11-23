//go:build manual

package wechat

import (
	"testing"
)

// 手动跑的。提前验证代码
func Test_service_manual_VerifyCode(t *testing.T) {
	//appId, ok := os.LookupEnv("WECHAT_APP_ID")
	//if !ok {
	//	panic("没有找到环境变量 WECHAT_APP_ID ")
	//}
	//appKey, ok := os.LookupEnv("WECHAT_APP_SECRET")
	//if !ok {
	//	panic("没有找到环境变量 WECHAT_APP_SECRET")
	//}
	//svc := NewService(appId, appKey)
	//res, err := svc.VerifyCode(context.Background(), "xxxx输入扫码登录后的 code", "state")
	//require.NoError(t, err)
	//t.Log(res)
}
