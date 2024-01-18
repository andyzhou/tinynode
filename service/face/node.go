package face

import (
	"errors"
	"github.com/andyzhou/tinylib/queue"
	"github.com/andyzhou/tinynode/define"
	"github.com/andyzhou/tinynode/json"
	"sync"
	"time"
)

/*
* node data face
 * @author <AndyZhou>
 * @mail <diudiu8848@163.com>
* - used for service side
 */

//face info
type Node struct {
	syncQueue *queue.Queue
	nodeMap sync.Map //remoteAddr -> *NodeInfo
}

//construct
func NewNode() *Node {
	this := &Node{
		syncQueue: queue.NewQueue(),
		nodeMap: sync.Map{},
	}
	this.interInit()
	return this
}

//quit
func (f *Node) Quit() {
	if f.syncQueue != nil {
		f.syncQueue.Quit()
	}
}

//sync node
func (f *Node) SyncNode(
	remoteAddr string,
	req *json.SyncNodeReq) error {
	//check
	if remoteAddr == "" || req == nil || req.Address == "" {
		return errors.New("invalid parameter")
	}
	nodeInfo, err := f.getNodeInfo(remoteAddr)
	if err != nil {
		return err
	}
	if nodeInfo == nil {
		return errors.New("remote addr not init")
	}

	//sync node info
	nodeInfo.Tag = req.Tag
	nodeInfo.Address = req.Address
	nodeInfo.Group = req.Group
	if req.Stat > define.NodeStatOfNone &&
		req.Stat <= define.NodeStatOfBusy {
		nodeInfo.Stat = req.Stat
	}
	nodeInfo.ActiveTime = time.Now().Unix()
	f.nodeMap.Store(remoteAddr, nodeInfo)
	return nil
}

//get remote address
func (f *Node) GetRemoteAddr() []string {
	result := make([]string, 0)
	sf := func(k, v interface{}) bool {
		remoteAddr, ok := k.(string)
		if ok && remoteAddr != "" {
			result = append(result, remoteAddr)
		}
		return true
	}
	f.nodeMap.Range(sf)
	return result
}

//remove node
//call this when remote node down
func (f *Node) RemoveNode(
	remoteAddr string) error {
	//check
	if remoteAddr == "" {
		return errors.New("invalid parameter")
	}
	//remove from run env
	f.nodeMap.Delete(remoteAddr)
	return nil
}

//get node info
func (f *Node) GetNode(
	remoteAddr string) (*json.NodeInfo, error) {
	//check
	if remoteAddr == "" {
		return nil, errors.New("invalid parameter")
	}
	nodeInfo, err := f.getNodeInfo(remoteAddr)
	if err != nil {
		return nil, err
	}
	return nodeInfo, nil
}

//init new node
//call this when remote node up
func (f *Node) InitNode(
	remoteAddr string) error {
	//check
	if remoteAddr == "" {
		return errors.New("invalid parameter")
	}
	nodeInfo, _ := f.getNodeInfo(remoteAddr)
	if nodeInfo != nil {
		return errors.New("node had init")
	}
	//init new
	nodeInfo = json.NewNodeInfo()
	nodeInfo.RemoteAddr = remoteAddr
	nodeInfo.Stat = define.NodeStatOfActive
	nodeInfo.ActiveTime = time.Now().Unix()
	f.nodeMap.Store(remoteAddr, nodeInfo)
	return nil
}

///////////////
//private func
///////////////

//get node info by remote addr
func (f *Node) getNodeInfo(
	remoteAddr string) (*json.NodeInfo, error) {
	//check
	if remoteAddr == "" {
		return nil, errors.New("invalid parameter")
	}
	v, ok := f.nodeMap.Load(remoteAddr)
	if !ok || v == nil {
		return nil, errors.New("can't get node")
	}
	nodeInfo, subOk := v.(*json.NodeInfo)
	if !subOk || nodeInfo == nil {
		return nil, nil
	}
	return nodeInfo, nil
}

//cb for node opt
func (f *Node) cbForNodeOpt(data interface{}) ([]byte, error) {
	return nil, nil
}

//inter init
func (f *Node) interInit() {
	f.syncQueue.SetCallback(f.cbForNodeOpt)
}
