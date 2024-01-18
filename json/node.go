package json

import "github.com/andyzhou/tinylib/util"

/*
 * node json info
 */

//node in env info
//also used for notify
type NodeInfo struct {
	Address    string `json:"address"` //unique
	RemoteAddr string `json:"remoteAddr"`
	Tag        string `json:"tag"`
	Group      string `json:"group"`
	Stat       int    `json:"stat"`
	ActiveTime int64  `json:"activeTime"`
	util.BaseJson
}

type NodesInfo struct {
	Total int         `json:"total"`
	Nodes []*NodeInfo `json:"nodes"`
	util.BaseJson
}

//sync node req
type SyncNodeReq struct {
	Address    string `json:"address"`
	Tag        string `json:"tag"`
	Group      string `json:"group"`
	util.BaseJson
}

//remove node req
type RemoveNodeReq struct {
	RemoteAddr string `json:"remoteAddr"`
	util.BaseJson
}

//gen opt resp
type GenOptResp struct {
	Err string `json:"err"`
	util.BaseJson
}

//construct
func NewNodeInfo() *NodeInfo {
	this := &NodeInfo{}
	return this
}
func NewNodesInfo() *NodesInfo {
	this := &NodesInfo{
		Nodes: []*NodeInfo{},
	}
	return this
}

func NewSyncNodeReq() *SyncNodeReq {
	this := &SyncNodeReq{}
	return this
}
func NewRemoveNodeReq() *RemoveNodeReq {
	this := &RemoveNodeReq{}
	return this
}
func NewGenOptResp() *GenOptResp {
	this := &GenOptResp{}
	return this
}