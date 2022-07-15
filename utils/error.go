package utils

import "errors"

var (
	ErrorServer    = errors.New("服务器出了点问题，重试一下？")
	ErrorIdPwd     = errors.New("账号或者密码错误,请重新输入")
	ErrorJwc       = errors.New("教务系统出了点问题,请重试")
	ErrorInput     = errors.New("参数有误")
	ErrorNoStudent = errors.New("查无此人")
)
