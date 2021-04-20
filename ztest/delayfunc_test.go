package ztest

import (
	"fmt"
	"server/ztimer"
	"testing"
)

/*
	针对delayFunc.go做单元测试
	主要测试延迟函数结构体是否正常使用
*/

func SayHello(message ...interface{}) {
	fmt.Println(message[0].(string), " ", message[1].(string))
}

func TestDelayfunc(t *testing.T) {
	df := ztimer.NewDelayFunc(SayHello, []interface{}{"hello", "world!"})
	fmt.Println("df.String() = ", df.String())
	df.Call()
}
