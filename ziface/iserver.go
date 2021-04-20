package ziface

//定义服务器接口
type IServer interface {
	//启动服务器
	Start()
	//停止服务器
	Stop()
	//开启业务服务
	Serve()
	//得到连接管理
	GetConnMgr() IConnMgr
	//设置该Server的连接创建时Hook函数
	SetOnConnStart(func(IConn))
	//设置该Server的连接断开时的Hook函数
	SetOnConnStop(func(IConn))
	//调用连接OnConnStart Hook函数
	CallOnConnStart(conn IConn)
	//调用连接OnConnStop Hook函数
	CallOnConnStop(conn IConn)
	//路由功能：给当前服务注册一个路由业务方法，共客户端连接处理使用
	AddRouter(msgID uint32, router IRouter)
}
