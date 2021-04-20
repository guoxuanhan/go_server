package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"server/utils"
	"server/ziface"
)

//封包拆包类实例，暂时不需要成员字段
type DataPack struct {
}

//封包拆包实例初始化
func NewDataPack() *DataPack {
	return &DataPack{}
}

//获取包头长度
func (dp *DataPack) GetHeadLen() uint32 {
	//ID (uint32) + DataLen(uint32) = 8
	return 8
}

//封包（压缩数据）
func (dp *DataPack) Pack(msg ziface.IMsg) ([]byte, error) {
	//创建一个存放bytes字节缓冲区
	dataBuff := bytes.NewBuffer([]byte{})

	//写dataLen
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	//写msgID
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}

	//写data数据
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

//拆包（解压数据）
func (dp *DataPack) UnPack(binaryData []byte) (ziface.IMsg, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	//只解压head的信息，得到dataLen和msgID
	msg := &Msg{}

	//读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Len); err != nil {
		return nil, err
	}

	//读msgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.ID); err != nil {
		return nil, err
	}

	//判断dataLen的长度是否超出我们允许的最大包长度
	if utils.GlobalObject.MaxPacketSize > 0 && msg.Len > utils.GlobalObject.MaxPacketSize {
		return nil, errors.New("too large msg data received")
	}

	//这里只需要把head的数据拆包出来就可以了，然后再通过head的长度，再从conn读取一次数据
	return msg, nil
}
