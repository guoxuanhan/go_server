package ziface

//IRequest接口：
//实际上是把客户端请求的连接信息和请求的数据包装到了Request里
type IRequest interface {
	//获取请求连接信息
	GetConn() IConn
	//获取请求的消息数据
	GetData() []byte
	//获取请求的消息ID
	GetMsgID() uint32
}
