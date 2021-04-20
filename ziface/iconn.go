package ziface

import "net"

//定义连接接口
type IConn interface {
	//启动连接，让当前连接开始工作
	Start()
	//停止连接，结束当前连接状态
	Stop()
	//从当前连接获取原始的socket
	GetTCPConn() *net.TCPConn
	//获取当前连接ID
	GetConnID() uint32
	//获取远程客户端地址信息
	RemoteAddr() net.Addr
	//直接将Msg数据发送给远程的TCP客户端（无缓冲）
	SendMsg(msgID uint32, data []byte) error
	//直接将Message数据发送给远程的TCP客户端(有缓冲)
	SendBuffMsg(msgID uint32, data []byte) error
	//设置连接属性
	SetProperty(key string, value interface{})
	//获取连接属性
	GetProperty(key string) (interface{}, error)
	//删除连接属性
	DelProperty(key string)
}
