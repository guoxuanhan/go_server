package znet

import (
	"fmt"
	"net"
	"server/utils"
	"server/ziface"
)

//IServer接口实现，定义一个Server服务类
type Server struct {
	//服务器的名称
	Name string
	//服务器版本号(tcp4 or other)
	IPVersion string
	//服务器绑定的ip地址
	IP string
	//服务器绑定的端口
	Port int
	//当前Server的消息管理模块，用来绑定MsgID和对应的处理方法
	msgHandler ziface.IMsgHandle
	//当前Server的连接管理器
	ConnMgr ziface.IConnMgr
	//该Server的连接创建时Hook函数
	OnConnStart func(conn ziface.IConn)
	//该Server的连接断开时的Hook函数
	OnConnStop func(conn ziface.IConn)
}

//创建一个服务器句柄
func NewServer() ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		msgHandler: NewMsgHandle(),
		ConnMgr:    NewConnMgr(),
	}
	return s
}

//实现ziface.IServer中全部接口方法

//开启网络环境
func (s *Server) Start() {
	fmt.Printf("[START] Server name: %s,listenner at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)

	//开启一个go去做服务端Linster业务
	go func() {
		//0.启动worker工作池机制
		s.msgHandler.StartWorkerPool()

		//1.获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}
		//2.监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}
		//已经监听成功
		fmt.Println("start Server  ", s.Name, " succ, now listenning...")

		//生成连接ID
		var connID uint32
		connID = 0

		//3.启动server网络连接业务
		for {
			//3.1阻塞等待客户端建立连接请求
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}
			fmt.Println("Get conn remote addr = ", conn.RemoteAddr().String())

			//3.2设置服务器最大连接控制，如果超过最大连接，则关闭此当前新连接
			if s.ConnMgr.Len() > utils.GlobalObject.MaxConn {
				conn.Close()
				continue
			}
			//3.3处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
			dealConn := NewConn(s, conn, connID, s.msgHandler)
			connID++

			//3.4启动当前链接的处理业务
			go dealConn.Start()
		}
	}()
}

//停止网络并清理
func (s *Server) Stop() {
	fmt.Println("[STOP] Server , name ", s.Name)

	//将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
	s.ConnMgr.ClearConn()
}

//运行服务器
func (s *Server) Serve() {
	s.Start()
	//TODO Server.Serve() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加

	//阻塞,否则主Go退出， listenner的go将会退出
	select {}
}

//得到连接管理
func (s *Server) GetConnMgr() ziface.IConnMgr {
	return s.ConnMgr
}

//设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(ziface.IConn)) {
	s.OnConnStart = hookFunc
}

//设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hookFunc func(ziface.IConn)) {
	s.OnConnStop = hookFunc
}

//调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn ziface.IConn) {
	if s.OnConnStart != nil {
		fmt.Println("---> CallOnConnStart....")
		s.OnConnStart(conn)
	}
}

//调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn ziface.IConn) {
	if s.OnConnStop != nil {
		fmt.Println("---> CallOnConnStop....")
		s.OnConnStop(conn)
	}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.msgHandler.AddRouter(msgID, router)
}
