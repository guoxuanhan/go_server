package ztimer

import (
	"fmt"
	"reflect"
	"server/zlog"
)

/*
	定义一个延迟调用函数
	延迟调用函数就是：时间定时器超时的时候，触发先注册好的回调函数
*/
type DelayFunc struct {
	//f: 延迟函数调用原型
	f func(...interface{})
	//args: 延迟调用函数传递的形参
	args []interface{}
}

/*
	创建一个延迟调用函数
*/
func NewDelayFunc(f func(v ...interface{}), args []interface{}) *DelayFunc {
	return &DelayFunc{
		f:    f,
		args: args,
	}
}

//打印当前延迟函数的信息，用于日志记录
func (df *DelayFunc) String() string {
	return fmt.Sprintf("{DelayFun:%s, args:%v}", reflect.TypeOf(df.f).Name(), df.args)
}

//执行延迟函数（如果执行失败，抛出异常）
func (df *DelayFunc) Call() {
	defer func() {
		/*
			内建函数recover允许程序管理恐慌过程中的Go程。
			在defer的函数中，执行recover调用会取回传至panic调用的错误值，
			恢复正常执行，停止恐慌过程。
			若recover在defer的函数之外被调用，它将不会停止恐慌过程序列。
			在此情况下，或当该Go程不在恐慌过程中时，或提供给panic的实参为nil时，
			recover就会返回nil。
		*/
		if err := recover(); err != nil {
			zlog.Error(df.String(), "Call err:", err)
		}
	}()

	//调用定时器超时函数
	df.f(df.args...)
}
