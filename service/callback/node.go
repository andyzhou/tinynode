package callback

import (
	"errors"
	"github.com/andyzhou/tinynode/define"
	"github.com/andyzhou/tinynode/json"
	"github.com/andyzhou/tinynode/service/face"
	"github.com/andyzhou/tinyrpc/proto"
)

/*
 * callback for node data
 * @author <AndyZhou>
 * @mail <diudiu8848@163.com>
 * - process access node data request
 */

//face info
type Node struct {
}

//construct
func NewNode() *Node {
	this := &Node{}
	return this
}

//decode node sync request
//decode query doc info request
func (f *Node) DecodeSync(
	remoteAddr string,
	in *proto.Packet) *proto.Packet {
	var (
		out proto.Packet
	)
	//check
	if remoteAddr == "" || in == nil || in.Data == nil {
		out.ErrCode = define.ErrCodeOfInvalidPara
		return &out
	}

	//decode origin data
	docObj := json.NewSyncNodeReq()
	err := docObj.Decode(in.Data, docObj)
	if err != nil {
		out.ErrCode = define.ErrCodeOfInvalidData
		out.ErrMsg = err.Error()
		return &out
	}

	//call bucket base func
	finalOut := f.docBaseOpt(remoteAddr, in.MessageId, docObj)
	return finalOut
}

////////////////
//private func
////////////////

//doc base opt
func (f *Node) docBaseOpt(
		remoteAddr string,
		messageId int32,
		docObj interface{},
	) *proto.Packet {
	var (
		outData []byte
		out     proto.Packet
		err     error
	)

	//call doc opt by message id
	switch messageId {
	case define.MessageIdOfSyncNode:
		{
			//sync node info
			outData, err = f.syncNodeOpt(remoteAddr, docObj)
			if err != nil {
				out.ErrCode = define.ErrCodeOfSaveData
				out.ErrMsg = err.Error()
				return &out
			}
		}
	}

	//format out
	out.MessageId = messageId
	out.ErrCode = define.ErrCodeOfSucceed
	out.Data = outData
	return &out
}

//sync node info opt
func (f *Node) syncNodeOpt(
	remoteAddr string,
	dataObj interface{}) ([]byte, error) {
	//detect input request
	req, ok := dataObj.(*json.SyncNodeReq)
	if !ok || req == nil {
		return nil, errors.New("data should be `SyncNodeReq`")
	}

	//call base func to save node
	nodeFace := face.GetInterFace().GetNode()
	err := nodeFace.SyncNode(remoteAddr, req)
	if err != nil {
		return nil, err
	}

	//format out byte data
	respObj := json.NewGenOptResp()
	outData, _ := respObj.Encode(respObj)
	return outData, nil
}