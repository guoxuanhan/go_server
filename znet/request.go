package znet

import "server/ziface"

type Request struct {
	//已经和客户端建立好的连接
	conn ziface.IConn
	//客户端请求的数据
	msg ziface.IMsg
}

//获取请求连接信息
func (r *Request) GetConn() ziface.IConn {
	return r.conn
}

//获取请求信息的数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

//获取请求的消息ID
func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgID()
}
