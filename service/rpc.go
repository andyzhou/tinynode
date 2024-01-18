package service

import (
	"errors"
	"github.com/andyzhou/tinynode/cmd"
	"github.com/andyzhou/tinynode/define"
	"github.com/andyzhou/tinynode/json"
	"github.com/andyzhou/tinynode/service/callback"
	"github.com/andyzhou/tinynode/service/face"
	"github.com/andyzhou/tinyrpc"
	"github.com/andyzhou/tinyrpc/proto"
)

/*
 * rpc service
 * - base on rpc protocol
 */

//service info
type RpcService struct {
	cmdConf *cmd.ServiceCfg //refer from outside
	rpcService *tinyrpc.Service //rpc service
}

//construct
func NewRpcService(cmdConf *cmd.ServiceCfg) *RpcService {
	this := &RpcService{
		cmdConf: cmdConf,
	}
	this.interInit()
	return this
}

//quit
func (f *RpcService) Quit() {
	if f.rpcService != nil {
		f.rpcService.Quit()
	}
}

//cb for general rpc requests
//calls from client node
func (f *RpcService) cbForRpcGenReq(
		remoteAddr string,
		in *proto.Packet,
	) (*proto.Packet, error) {
	//check
	if remoteAddr == "" || in == nil {
		return nil, errors.New("invalid parameter")
	}

	//run gen rpc callback
	cb := callback.GetCallBack()
	out := cb.DecodePackage(remoteAddr, in)

	//special opt
	switch in.MessageId {
	case define.MessageIdOfSyncNode:
		{
			//notify all remote clients
			f.notifyNode(remoteAddr)
		}
	}
	return out, nil
}

//notify node to clients
//opt this by stream way
func (f *RpcService) notifyNode(
	remoteAddr string,
	nodeInfos ...*json.NodeInfo) error {
	var (
		nodeInfo *json.NodeInfo
	)
	//check
	if remoteAddr == "" {
		return errors.New("invalid parameter")
	}
	if nodeInfos != nil && len(nodeInfos) > 0 {
		nodeInfo = nodeInfos[0]
	}

	//get node info by remote addr
	nodeFace := face.GetInterFace().GetNode()
	if nodeInfo == nil {
		nodeInfo, _ = nodeFace.GetNode(remoteAddr)
	}
	if nodeInfo == nil {
		return errors.New("can't get node info by remote addr")
	}

	//setup pack
	objEnc, _ := nodeInfo.Encode(nodeInfo)
	packObj := &proto.Packet{
		MessageId: define.MessageIdOfNotifyNode,
		Data: objEnc,
	}

	//cast to remote
	err := f.rpcService.SendStreamData(packObj)
	return err
}

//cb for client node down
func (f *RpcService) cbForClientNodeDown(
	remoteAddr string) error {
	//remove down node
	nodeFace := face.GetInterFace().GetNode()
	nodeObj, _ := nodeFace.GetNode(remoteAddr)
	if nodeObj != nil {
		//notify to client
		nodeObj.Stat = define.NodeStatOfDown
		f.notifyNode(remoteAddr, nodeObj)
	}
	//remove node
	err := nodeFace.RemoveNode(remoteAddr)
	return err
}

//cb for client node up
func (f *RpcService) cbForClientNodeUp(
	remoteAddr string) error {
	//init new node
	nodeFace := face.GetInterFace().GetNode()
	err := nodeFace.InitNode(remoteAddr)
	return err
}

//register rpc service
func (f *RpcService) registerRpc() {
	//get rpc port
	rpcPort := f.cmdConf.RpcPort

	//init rpc service
	f.rpcService = tinyrpc.NewService()

	//set relate callback
	f.rpcService.SetCBForClientNodeUp(f.cbForClientNodeUp)
	f.rpcService.SetCBForClientNodeDown(f.cbForClientNodeDown)
	f.rpcService.SetCBForGeneral(f.cbForRpcGenReq)

	//start rpc service
	err := f.rpcService.Start(rpcPort)
	if err != nil {
		panic(any(err))
	}
}

//inter init
func (f *RpcService) interInit() {
	//register rpc service
	f.registerRpc()
}