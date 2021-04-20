package znet

import (
	"errors"
	"fmt"
	"server/ziface"
	"sync"
)

//连接管理模块
type ConnMgr struct {
	//管理的连接信息
	conns map[uint32]ziface.IConn
	//读写连接的读写锁
	connsLock sync.RWMutex
}

//创建一个连接管理
func NewConnMgr() *ConnMgr {
	return &ConnMgr{
		conns: make(map[uint32]ziface.IConn),
	}
}

//添加连接
func (cm *ConnMgr) Add(conn ziface.IConn) {
	cm.connsLock.Lock()
	defer cm.connsLock.Unlock()

	//将conn添加到cm.conns中管理
	cm.conns[conn.GetConnID()] = conn
	fmt.Println("connection add to ConnManager successfully: conn num = ", cm.Len())
}

//删除连接
func (cm *ConnMgr) Del(conn ziface.IConn) {
	cm.connsLock.Lock()
	defer cm.connsLock.Unlock()

	//删除连接信息
	delete(cm.conns, conn.GetConnID())
	fmt.Println("connection Remove ConnID=", conn.GetConnID(), " successfully: conn num = ", cm.Len())
}

//通过connID获取连接
func (cm *ConnMgr) Get(connID uint32) (ziface.IConn, error) {
	//保护共享资源Map 加读锁
	cm.connsLock.RLock()
	defer cm.connsLock.RUnlock()

	if conn, ok := cm.conns[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

//获取当前连接个数
func (cm *ConnMgr) Len() int {
	return len(cm.conns)
}

//删除并停止所有连接
func (cm *ConnMgr) ClearConn() {
	cm.connsLock.Lock()
	defer cm.connsLock.Unlock()

	for connID, conn := range cm.conns {
		//停止
		conn.Stop()
		//删除
		delete(cm.conns, connID)
	}
	fmt.Println("Clear All Connections successfully: conn num = ", cm.Len())
}
