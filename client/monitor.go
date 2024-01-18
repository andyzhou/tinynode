package client

import (
	"errors"
	"github.com/andyzhou/tinynode/define"
	"github.com/andyzhou/tinynode/json"
	"github.com/andyzhou/tinyrpc/proto"
	"sync"
)

/*
 * monitor face
 * - connect one monitor and sync data
 * - cache the nodes of one monitor server
 * - other opt with monitor server
 */

//face info
type Monitor struct {
	monitorAddr string //monitor address
	clientNodeMap map[string]*Node //nodeAddr -> *Node
	cbForNodeNotify func(*json.NodeInfo) error
	sync.RWMutex
}

//construct
func NewMonitor(addr string) *Monitor {
	this := &Monitor{
		monitorAddr: addr,
		clientNodeMap: map[string]*Node{},
	}
	return this
}

//quit
func (f *Monitor) Quit() {
	f.Lock()
	defer f.Unlock()
	for _, v := range f.clientNodeMap {
		v.Quit()
	}
	f.clientNodeMap = map[string]*Node{}
}

//get batch nodes info
//support filter
func (f *Monitor) GetNodesInfo(
		groups ...string,
	) ([]*json.NodeInfo, error) {
	var (
		group string
	)
	if groups != nil && len(groups) > 0 {
		group = groups[0]
	}

	//format result
	result := make([]*json.NodeInfo, 0)

	//loop filter with locker
	f.Lock()
	defer f.Unlock()
	for _, v := range f.clientNodeMap {
		nodeInfo := v.GetNodeInfo()
		if group != "" {
			//filter by group
			if nodeInfo.Group == group {
				result = append(result, nodeInfo)
			}
		}else{
			//general
			result = append(result, nodeInfo)
		}
	}
	return result, nil
}

//sync client node info to monitor
//this call from client side manually
func (f *Monitor) SyncNodeInfo(
	req *json.SyncNodeReq) error {
	//check
	if req == nil || req.Address == "" {
		return errors.New("invalid parameter")
	}

	//get and init node by addr
	node, err := f.getNodeByAddr(req.Address)
	if err != nil {
		return err
	}
	if node == nil {
		//init new node
		node, err = f.initNewNode(req)
		if err != nil {
			return err
		}
	}

	//gen packet
	objByte, _ := req.Encode(req)
	pack := &proto.Packet{
		MessageId: define.MessageIdOfSyncNode,
		Data: objByte,
	}

	//sync node info to monitor server
	_, err = node.senGenRpcRequest(pack)
	return err
}

//set cb for node notify
func (f *Monitor) SetCBForNodeNotify(cb func(info *json.NodeInfo) error) {
	if cb == nil {
		return
	}
	f.cbForNodeNotify = cb
}

////////////////
//private func
////////////////

//inter cb for node notify
func (f *Monitor) interCBForNodeNotify(info *json.NodeInfo) error {
	//check
	if info == nil || info.Address == "" {
		return errors.New("invalid parameter")
	}

	//sync into run env data
	nodeObj, _ := f.getNodeByAddr(info.Address)
	if nodeObj != nil {
		nodeObj.SyncNodeInfo(info)
	}

	//check and call cb of outside
	if f.cbForNodeNotify != nil {
		f.cbForNodeNotify(info)
	}
	return nil
}

//init new node
func (f *Monitor) initNewNode(
	req *json.SyncNodeReq) (*Node, error) {
	//check
	if req == nil || req.Address == "" {
		return nil, errors.New("invalid parameter")
	}
	//init new node
	node := NewNode()
	err := node.InitNode(f.monitorAddr, req, f.interCBForNodeNotify)
	if err != nil {
		return nil, err
	}
	//sync into run env map
	f.Lock()
	defer f.Unlock()
	f.clientNodeMap[req.Address] = node
	return node, nil
}

//get node by address
func (f *Monitor) getNodeByAddr(
	addr string) (*Node, error) {
	//check
	if addr == "" {
		return nil, errors.New("invalid parameter")
	}

	//get node by address
	f.Lock()
	defer f.Unlock()
	v, ok := f.clientNodeMap[addr]
	if ok && v != nil {
		return v, nil
	}
	return v, nil
}
