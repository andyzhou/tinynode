package callback

import (
	"github.com/andyzhou/tinynode/define"
	"github.com/andyzhou/tinyrpc/proto"
	"sync"
)

/*
 * rpc service callback face
 * - main entry face
 */

//global variable
var (
	_cb *CallBack
	_cbOnce sync.Once
)

//face info
type CallBack struct {
	node *Node
}

//get single instance
func GetCallBack() *CallBack {
	_cbOnce.Do(func() {
		_cb = NewCallBack()
	})
	return _cb
}

//construct
func NewCallBack() *CallBack {
	//self init
	this := &CallBack{
		node: NewNode(),
	}
	return this
}

//decode and process request package
//return *proto.Packet, error
func (f *CallBack) DecodePackage(
	remoteAddr string,
	in *proto.Packet) *proto.Packet {
	var (
		out = &proto.Packet{}
	)
	//check
	if remoteAddr == "" || in == nil ||
		in.MessageId <= define.MessageIdOfNone {
		out.ErrCode = define.ErrCodeOfInvalidPara
		return out
	}

	//setup out value
	out.MessageId = in.MessageId

	//call sub func by message id
	switch in.MessageId {
	case define.MessageIdOfSyncNode:
		{
			//sync node info
			out = f.node.DecodeSync(remoteAddr, in)
		}
	}

	return out
}