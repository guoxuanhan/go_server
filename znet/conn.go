package znet

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"server/utils"
	"server/ziface"
	"sync"
)

type Conn struct {
	//当前Conn属于哪个Server
	TcpServer ziface.IServer
	//当前连接的socket套接字
	Conn *net.TCPConn
	//当前连接的ID（也可以称作为seccionID，iD全局唯一）
	ConnID uint32
	//告知该链接已经退出/停止的channel
	ctx    context.Context
	cancel context.CancelFunc
	//当前连接的关闭状态
	isClosed bool
	//消息管理MsgID和对应的处理方法的消息管理模块
	MsgHandler ziface.IMsgHandle
	//无缓冲管道，用于读写两个goroutine之间的消息通信
	msgChan chan []byte
	//有缓冲管道，用于读、写两个goroutine之间的消息通信
	msgBuffChan chan []byte
	//读写锁
	sync.RWMutex
	//连接属性
	property map[string]interface{}
	//保护当前property的锁
	propertyLock sync.Mutex
}

//创建连接的方法
func NewConn(server ziface.IServer, conn *net.TCPConn, connID uint32, msghandler ziface.IMsgHandle) *Conn {
	//初始化Conn属性
	c := &Conn{
		TcpServer:   server,
		Conn:        conn,
		ConnID:      connID,
		isClosed:    false,
		MsgHandler:  msghandler,
		msgChan:     make(chan []byte),
		msgBuffChan: make(chan []byte, utils.GlobalObject.MaxMsgChanLen),
		property:    make(map[string]interface{}),
	}

	//将新创建的conn添加到连接管理器中
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

//写消息的goroutine，用户将数据发送给客户端
func (c *Conn) Writer() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")

	for {
		select {
		case data := <-c.msgChan:
			//有（无缓冲）数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send Data error:, ", err, " Conn Writer exit")
				return
			}
		case data, ok := <-c.msgBuffChan:
			if ok {
				//有数据要写给客户端
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("Send Buff Data error:, ", err, " Conn Writer exit")
					return
				}
			} else {
				fmt.Println("msgBuffChan is Closed")
				break
			}
		case <-c.ctx.Done():
			return
		}
	}
}

//读消息的goroutine，用于从客户端读取数据
func (c *Conn) Reader() {
	fmt.Println("[Reader Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Reader exit!]")
	defer c.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			//创建拆包解包的对象
			dp := NewDataPack()
			//读取客户端的MsgHead
			headData := make([]byte, dp.GetHeadLen())
			if _, err := io.ReadFull(c.Conn, headData); err != nil {
				fmt.Println("read msg head error ", err)
				return
			}
			//拆包得到msgID和dataLen 放在msg中
			msg, err := dp.UnPack(headData)
			if err != nil {
				fmt.Println("unpack error ", err)
				return
			}
			//根据 dataLen 读取 data，放在msg.Data中
			var data []byte
			if msg.GetDataLen() > 0 {
				data = make([]byte, msg.GetDataLen())
				if _, err := io.ReadFull(c.Conn, data); err != nil {
					fmt.Println("read msg data error ", err)
					return
				}
			}
			msg.SetData(data)
			//得到当前客户端请求的Request数据
			req := Request{
				conn: c,
				msg:  msg,
			}
			if utils.GlobalObject.WorkerPoolSize > 0 {
				//已经启动工作池机制，将消息交给Worker处理（阻塞排队执行）
				c.MsgHandler.SendMsgToTaskQueue(&req)
			} else {
				//从绑定好的消息和对应的处理方法中执行对应的Handle方法（有请求立即执行）
				go c.MsgHandler.DoMsgHandler(&req)
			}
		}
	}
}

//启动连接，让当前连接开始工作
func (c *Conn) Start() {
	c.ctx, c.cancel = context.WithCancel(context.Background())
	//1.开启用户从客户端读取数据的goroutine
	go c.Reader()
	//2.开启用于写回客户端数据流程的goroutine
	go c.Writer()
	//按照用户传递进来的创建连接时需要处理的业务，执行Hook方法
	c.TcpServer.CallOnConnStart(c)
}

//停止连接，结束当前连接状态
func (c *Conn) Stop() {
	fmt.Println("Conn Stop()...ConnID = ", c.ConnID)
	c.Lock()
	defer c.Unlock()

	//如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用
	c.TcpServer.CallOnConnStop(c)

	//如果当前链接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	//关闭socket链接
	c.Conn.Close()
	//关闭writer
	c.cancel()
	//将该连接从连接管理器中删除
	c.TcpServer.GetConnMgr().Del(c)
	//关闭该连接全部管道
	close(c.msgBuffChan)
	close(c.msgChan)
}

//从当前连接获取原始的socket
func (c *Conn) GetTCPConn() *net.TCPConn {
	return c.Conn
}

//获取当前连接ID
func (c *Conn) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端地址信息
func (c *Conn) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//直接将Msg数据发送给远程的TCP客户端（无缓冲）
func (c *Conn) SendMsg(msgID uint32, data []byte) error {
	c.RLock()
	defer c.RUnlock()
	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}
	//将data封包并发送
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgID, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgID)
		return errors.New("Pack error msg ")
	}
	//写回客户端
	c.msgChan <- msg
	return nil
}

//直接将Message数据发送给远程的TCP客户端(有缓冲)
func (c *Conn) SendBuffMsg(msgID uint32, data []byte) error {
	c.RLock()
	defer c.RUnlock()
	if c.isClosed == true {
		return errors.New("Connection closed when send buff msg")
	}
	//将data封包并发送
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgID, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgID)
		return errors.New("Pack error msg ")
	}
	//写回客户端
	c.msgBuffChan <- msg
	return nil
}

//设置连接属性
func (c *Conn) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

//获取连接属性
func (c *Conn) GetProperty(key string) (interface{}, error) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

//删除连接属性
func (c *Conn) DelProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
