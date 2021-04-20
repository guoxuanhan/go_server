package znet

type Msg struct {
	//消息的长度
	Len uint32
	//消息的ID
	ID uint32
	//消息的内容
	Data []byte
}

//创建一个Msg消息包
func NewMsgPackage(id uint32, data []byte) *Msg {
	return &Msg{
		Len:  uint32(len(data)),
		ID:   id,
		Data: data,
	}
}

//获取消息数据段长度
func (msg *Msg) GetDataLen() uint32 {
	return msg.Len
}

//获取消息ID
func (msg *Msg) GetMsgID() uint32 {
	return msg.ID
}

//获取消息内容
func (msg *Msg) GetData() []byte {
	return msg.Data
}

//设置消息ID
func (msg *Msg) SetMsgID(id uint32) {
	msg.ID = id
}

//设置消息内容
func (msg *Msg) SetData(data []byte) {
	msg.Data = data
}

//设置消息数据段长度
func (msg *Msg) SetDataLen(len uint32) {
	msg.Len = len
}
