package ziface

//连接管理抽象层
type IConnMgr interface {
	//添加连接
	Add(conn IConn)
	//删除连接
	Del(conn IConn)
	//通过connID获取连接
	Get(connID uint32) (IConn, error)
	//获取当前连接个数
	Len() int
	//删除并停止所有连接
	ClearConn()
}
