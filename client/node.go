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

//sync node info
//send to server and store in run env
func (f *Node) SyncNodeInfo(
	req *json.SyncNodeReq) error {
	//check
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
	serverAddr string) error {
	//check
	if serverAddr == "" {
		return errors.New("invalid parameter")
	}

	//init new rpc client
	newClient := tinyrpc.NewClient()
	err := newClient.SetAddress(serverAddr)
	if err != nil {
		return err
	}

	//set callback and connect server
	newClient.SetServerNodeDownCallBack(f.cbForServerNodeDown)
	err = newClient.ConnectServer()
	if err != nil {
		return err
	}

	//sync rpc client
	f.Lock()
	defer f.Unlock()
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

	////update node stat
	//if f.info != nil {
	//	f.Lock()
	//	f.info.Stat = define.NodeStatOfDown
	//	f.Unlock()
	//}

	//re-connect target server force
	//run in loop?
	err := f.client.ConnectServer()
	if err != nil {
		log.Printf("node.CBForServerNodeDown, serverAddr:%v, err:%v\n", serverAddr, err.Error())
	}
	return err
}
