package util

import "github.com/gin-gonic/gin"

// R 统一 API 响应结构，等价 Java 的 R<T>
type R struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func Ok(c *gin.Context, data any) {
	c.JSON(200, R{Code: 200, Msg: "ok", Data: data})
}

func OkEmpty(c *gin.Context) {
	c.JSON(200, R{Code: 200, Msg: "ok", Data: nil})
}

func Fail(c *gin.Context, msg string) {
	c.JSON(200, R{Code: 500, Msg: msg, Data: nil})
}

func FailCode(c *gin.Context, code int, msg string) {
	c.JSON(200, R{Code: code, Msg: msg, Data: nil})
}

func Unauthorized(c *gin.Context) {
	c.JSON(200, R{Code: 401, Msg: "未登录或Token已过期", Data: nil})
}
