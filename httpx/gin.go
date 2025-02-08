package httpx

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 响应 结构体
type Reply struct {
	Code   int         `json:"Code"`             // 结果码: 正常=0 失败>0
	Status string      `json:"Status,omitempty"` // 正常=OK，失败=具体错误内容
	Error  string      `json:"Error,omitempty"`  // 错误信息
	Msg    string      `json:"Msg,omitempty"`    // 信息提示
	RunMS  int64       `json:"RunMS,omitempty"`  // 执行时间，单位毫秒 ms
	Data   interface{} `json:"Data,omitempty"`   // 结果值
	Err    error       `json:"-"`                // 错误
}

func NewReply() Reply {
	return Reply{Code: 0, Status: "OK"} // Code: 0 时正确响应了请求，成功
}

func Return(msg string, reply *Reply, start *time.Time) {
	reply.Msg = fmt.Sprintf("%s 成功", msg)
	reply.Code = 0
	reply.RunMS = time.Since(*start).Milliseconds()
	if reply.Err != nil {
		reply.Code = 400
		reply.Status = "ERROR"
		reply.Error = reply.Err.Error()

		if _, file, line, ok := runtime.Caller(2); ok {
			// reply.Error = fmt.Sprintf("[error] %s:%s:%d %+v", runtime.FuncForPC(pc).Name(), file, line, reply.Err)
			// reply.Error = fmt.Sprintf("%s [error] %s:%s:%d", reply.Err.Error(), runtime.FuncForPC(pc).Name(), file, line)
			reply.Error = fmt.Sprintf("%s [error] %s:%d", reply.Err.Error(), file, line)
		}
		reply.Msg = fmt.Sprintf("%s 失败", msg)
		reply.Data = nil
	}
}

func Recover() {
	if r := recover(); r != nil {
		log.Panicf("%s\r\n", strings.Repeat("!", 30))
		log.Panicf("SYSTEM ACTION PANIC: %v, STACK: %v", r, string(debug.Stack()))
		log.Panicf("%s\r\n", strings.Repeat("!", 30))
	}
}

func ReturnGin(c *gin.Context, msg string, reply *Reply, start *time.Time) {
	reply.Msg = fmt.Sprintf("%s 成功", msg)
	reply.Code = 0
	if start != nil {
		reply.RunMS = time.Since(*start).Milliseconds()
	}
	if reply.Err != nil {
		reply.Code = 400
		reply.Status = "ERROR"
		reply.Error = reply.Err.Error()
		if _, file, line, ok := runtime.Caller(2); ok {
			reply.Error = fmt.Sprintf("%s [error] %s:%d", reply.Err.Error(), file, line)
		}
		reply.Msg = fmt.Sprintf("%s 失败", msg)
		reply.Data = nil
	}
	c.JSON(http.StatusOK, reply)
}

func RecoverGin(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			errorMsg := fmt.Sprintf("panic: %v\n%s", r, stack)
			fmt.Printf("panic: %v\n%s", r, stack)
			c.AbortWithStatusJSON(400,
				gin.H{
					"Msg":    "错误",
					"Code":   400,
					"Status": "ERROR",
					"Error":  errorMsg,
				},
			)
		}
	}()
	c.Next()
}

func CorsGin() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token,X-Token,X-User-Id")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type, New-Token, New-Expires-At")
		c.Header("Access-Control-Allow-Credentials", "true")
		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}
