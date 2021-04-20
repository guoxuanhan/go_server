package ziface

//将请求的一个消息封装到msg中，定义抽象层接口
type IMsg interface {
	//获取消息数据段长度
	GetDataLen() uint32
	//获取消息ID
	GetMsgID() uint32
	//获取消息内容
	GetData() []byte
	//设置消息ID
	SetMsgID(uint32)
	//设置消息内容
	SetData([]byte)
	//设置消息数据段长度
	SetDataLen(uint32)
}
