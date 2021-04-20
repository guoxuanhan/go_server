package ziface

//路由接口：这里面路由是使用框架着给该连接自定义的处理业务方法
//路由里的IRequest则包含该连接的连接信息和该连接的请求数据信息
type IRouter interface {
	//在处理conn业务之前的钩子方法
	PreHandle(request IRequest)
	//处理conn业务的方法
	Handle(request IRequest)
	//处理conn业务之后的钩子方法
	PostHandle(request IRequest)
}
