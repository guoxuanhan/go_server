package znet

import (
	"fmt"
	"server/utils"
	"server/ziface"
	"strconv"
)

type MsgHandle struct {
	//存放每个MsgID所对应所对应的处理方法的map属性
	APIS map[uint32]ziface.IRouter
	//处理业务工作Worker池的数量
	WorkerPoolSize uint32
	//Worker负责取任务的消息队列
	TaskQueue []chan ziface.IRequest
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		APIS:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		//一个worker对应一个queue
		TaskQueue: make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

//马上以非阻塞方式处理消息
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	handler, ok := mh.APIS[request.GetMsgID()]
	if !ok {
		fmt.Println("APIS msgId = ", request.GetMsgID(), " is not FOUND!")
		return
	}
	//执行对应处理方法
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

//为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	//1.判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.APIS[msgID]; ok {
		panic("repeated api , msgId = " + strconv.Itoa(int(msgID)))
	}
	//2.添加msg与API的绑定关系
	mh.APIS[msgID] = router
	fmt.Println("Add API msgId = ", msgID)
}

//启动worker工作池
func (mh *MsgHandle) StartWorkerPool() {
	//遍历需要启动worker的数量，依此启动
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//给当前worker对应的任务队列开辟空间
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//启动当前worker，阻塞的等待对应的任务队列是否有消息传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

//将消息交给TaskQueue,由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	//根据ConnID来分配当前的连接应该由哪个worker负责处理
	//轮询的平均分配法则（求余）

	//得到需要处理此条连接的workerID
	workerID := request.GetConn().GetConnID() % mh.WorkerPoolSize
	//将请求消息发送给人物队列
	mh.TaskQueue[workerID] <- request
}

func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workerID, " is started.")
	//不断的等待队列中的消息
	for {
		select {
		//有消息则取出队列的Request，并执行绑定的业务方法
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}
