package zlog

import "os"

//全局默认提供一个log对外句柄，可以直接使用API系列调用
//全局日志对象 StdLog
var StdLog = NewLogger(os.Stderr, "", BitDefault)

//获取StdLog标记位
func Flags() int {
	return StdLog.Flags()
}

//设置StdLog标记位
func ResetFlags(flag int) {
	StdLog.ResetFlags(flag)
}

//添加flag标记
func AddFlag(flag int) {
	StdLog.AddFlag(flag)
}

//设置StdLog 日志头前缀
func SetPrefix(prefix string) {
	StdLog.SetPrefix(prefix)
}

//设置StdLog绑定的日志文件
func SetLogFile(fileDir string, filename string) {
	StdLog.SetLogFile(fileDir, filename)
}

//设置关闭debug
func CloseDebug() {
	StdLog.CloseDebug()
}

//设置打开debug
func OpenDebug() {
	StdLog.OpenDebug()
}

// ====> Debug <====
func Debugf(format string, v ...interface{}) {
	StdLog.Debugf(format, v...)
}

func Debug(v ...interface{}) {
	StdLog.Debug(v...)
}

// ====> Info <====
func Infof(format string, v ...interface{}) {
	StdLog.Infof(format, v...)
}

func Info(v ...interface{}) {
	StdLog.Info(v...)
}

// ====> Warn <====
func Warnf(format string, v ...interface{}) {
	StdLog.Warnf(format, v...)
}

func Warn(v ...interface{}) {
	StdLog.Warn(v...)
}

// ====> Error <====
func Errorf(format string, v ...interface{}) {
	StdLog.Errorf(format, v...)
}

func Error(v ...interface{}) {
	StdLog.Error(v...)
}

// ====> Panic  <====
func Panicf(format string, v ...interface{}) {
	StdLog.Panicf(format, v...)
}

func Panic(v ...interface{}) {
	StdLog.Panic(v...)
}

// ====> Fatal 需要终止程序 <====
func Fatalf(format string, v ...interface{}) {
	StdLog.Fatalf(format, v...)
}

func Fatal(v ...interface{}) {
	StdLog.Fatal(v...)
}

// ====> Stack  <====
func Stack(v ...interface{}) {
	StdLog.Stack(v...)
}

func init() {
	//因为StdLog对象 对所有输出方法做了一层包裹，所以在打印调用函数的时候，比正常的logger对象多一层调用
	//所以一般的Logger对象 calledDepth=2, StdLog的calldDepth=3
	StdLog.calledDepth = 3
}
