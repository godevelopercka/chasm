package web

// 501001 => 这里代表验证码错误, 5系统错误 01代表user模块, 001代表验证码
type Result struct {
	// 这个叫做业务错误码
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}
