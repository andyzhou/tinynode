package client

import (
	genJson "encoding/json"
	"errors"
	"github.com/andyzhou/tinynode/define"
	"github.com/andyzhou/tinynode/json"
	"github.com/andyzhou/tinyrpc"
	"github.com/andyzhou/tinyrpc/proto"
	"log"
	"sync"
	"time"
)

/*
 * node face
 * - one node opt and data face
 */

//face info
type Node struct {
	info *json.NodeInfo //node info
	client *tinyrpc.Client //rpc client obj
	serverDown bool
	cbForNodeNotify func(info *json.NodeInfo) error //reference
	sync.RWMutex
}

//construct
func NewNode() *Node {
	this := &Node{}
	return this
}

//quit
func (f *Node) Quit() {
	if f.client != nil {
		f.client.Quit()
	}
}

//get node client
func (f *Node) GetNodeClient() *tinyrpc.Client {
	return f.client
}

//get node info
func (f *Node) GetNodeInfo() *json.NodeInfo {
	return f.info
}

//sync node info into run env
func (f *Node) SyncNodeInfo(info *json.NodeInfo) error {
	//check
	if info == nil || info.Address == "" {
		return errors.New("invalid parameter")
	}
	if f.info == nil {
		return errors.New("node info in run env is nil")
	}
	//sync node info
	f.Lock()
	defer f.Unlock()
	f.info.Tag = info.Tag
	f.info.Group = info.Group
	f.info.Stat = info.Stat
	f.info.ActiveTime = time.Now().Unix()
	return nil
}

//send stream rpc data
func (f *Node) SendStreamRequest(
	messageId int32,
	jsonObj interface{}) error {
	//check
	if messageId <= define.MessageIdOfNone || jsonObj == nil {
		return errors.New("invalid parameter")
	}
	if f.client == nil {
		return errors.New("client not init")
	}

	//gen packet
	objByte, _ := genJson.Marshal(jsonObj)
	pack := &proto.Packet{
		MessageId: messageId,
		Data: objByte,
	}

	//begin send stream data to server
	err := f.client.SendStreamData(pack)
	return err
}

//send gen rpc request
func (f *Node) SendGenRequest(
		messageId int32,
		jsonObj interface{},
	) (*proto.Packet, error) {
	//check
	if messageId <= define.MessageIdOfNone || jsonObj == nil {
		return nil, errors.New("invalid parameter")
	}

	//gen packet
	objByte, _ := genJson.Marshal(jsonObj)
	pack := &proto.Packet{
		MessageId: messageId,
		Data: objByte,
	}

	//send real rpc request
	resp, err := f.senGenRpcRequest(pack)
	return resp, err
}

//init new node to connect server
func (f *Node) InitNode(
		serverAddr string,
		req *json.SyncNodeReq,
		cbForNodeNotify func(info *json.NodeInfo) error,
	) error {
	//check
	if serverAddr == "" || req == nil {
		return errors.New("invalid parameter")
	}

	//init new rpc client
	newClient := tinyrpc.NewClient()
	err := newClient.SetAddress(serverAddr)
	if err != nil {
		return err
	}
	if cbForNodeNotify != nil {
		f.cbForNodeNotify = cbForNodeNotify
	}

	//set callback
	newClient.SetServerNodeDownCallBack(f.cbForServerNodeDown)
	newClient.SetStreamCallBack(f.cbForReceiveStreamData)

	//connect server
	err = newClient.ConnectServer()
	if err != nil {
		return err
	}

	//init node info
	nodeInfo := json.NewNodeInfo()
	nodeInfo.Address = req.Address
	nodeInfo.Tag = req.Tag
	nodeInfo.Group = req.Group
	nodeInfo.Stat = define.NodeStatOfActive
	nodeInfo.ActiveTime = time.Now().Unix()
	if req.Stat > define.NodeStatOfNone {
		nodeInfo.Stat = req.Stat
	}

	//sync node info and rpc client
	f.Lock()
	defer f.Unlock()
	f.info = nodeInfo
	f.client = newClient
	return nil
}

///////////////
//private func
///////////////

//send gen rpc call
func (f *Node) senGenRpcRequest(
		pack *proto.Packet,
	) (*proto.Packet, error) {
	//check
	if pack == nil {
		return nil, errors.New("invalid parameter")
	}
	if f.client == nil {
		return nil, errors.New("client not init")
	}

	//send real request to target node
	resp, subErrThree := f.client.SendRequest(pack)
	return resp, subErrThree
}

//cb for receive stream data
func (f *Node) cbForReceiveStreamData(
	in *proto.Packet) error {
	var (
		err error
	)
	//check
	if in == nil || in.MessageId <= define.MessageIdOfNone ||
		in.Data == nil {
		return errors.New("invalid parameter")
	}

	//decode data
	nodeObj := json.NewNodeInfo()
	nodeObj.Decode(in.Data, nodeObj)
	if nodeObj == nil || nodeObj.Address == "" {
		return errors.New("invalid pack data")
	}

	//check and run cb for node notify
	if f.cbForNodeNotify != nil {
		err = f.cbForNodeNotify(nodeObj)
	}
	return err
}

//cb for server node down
func (f *Node) cbForServerNodeDown(
	serverAddr string) error {
	//check
	if serverAddr == "" {
		return errors.New("invalid parameter")
	}

	//update inter stat
	f.Lock()
	f.serverDown = true
	f.Unlock()

	//re-connect target server force
	//run in loop?
	err := f.client.ConnectServer()
	if err != nil {
		log.Printf("node.CBForServerNodeDown, serverAddr:%v, err:%v\n", serverAddr, err.Error())
	}
	return err
}
