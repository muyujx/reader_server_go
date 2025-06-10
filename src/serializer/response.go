package serializer

import (
	"github.com/gin-gonic/gin"
	"muyu.com/reader_server_go/v1/src/logger"
	"net/http"
)

type Response struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

var (
	Success = Response{
		Code: 0,
	}

	// UnknownErr 未知错误, 通常为未定义的内部错误
	UnknownErr = Response{
		Code: 1,
		Msg:  "unknown err",
	}

	// NeedLogin 需要登录
	NeedLogin = Response{
		Code: 100,
		Msg:  "need login",
	}

	// ParamErr 请求参数解析失败
	ParamErr = Response{
		Code: 101,
		Msg:  "param err",
	}

	// AccountExits 注册时账号已经存在
	AccountExits = Response{
		Code: 102,
		Msg:  "account already exits",
	}

	// NoPermission 没有权限
	NoPermission = Response{
		Code: 103,
		Msg:  "no permission",
	}
)

func ReturnJson(c *gin.Context, data Response) {
	c.JSON(http.StatusOK, data)
}

func SuccessRes(obj any) Response {
	return Response{
		Code: 0,
		Data: obj,
	}
}

func SuccessNoRes() Response {
	return Response{
		Code: 0,
		Data: nil,
	}
}

func ErrRes(err error) Response {
	logger.Error("error", err)
	return Response{
		Code: 1,
		Msg:  "error",
	}
}

func ErrMsgRes(msg string) Response {
	return Response{
		Code: -1,
		Msg:  msg,
	}
}

func ParamErrMsg(msg string) Response {
	return Response{
		Code: 101,
		Msg:  msg,
	}
}

func ErrCodeMsgRes(code int, msg string) Response {
	return Response{
		Code: code,
		Msg:  msg,
	}
}
