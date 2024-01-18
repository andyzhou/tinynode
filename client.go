package tinynode

import (
	"errors"
	"github.com/andyzhou/tinynode/client"
	"github.com/andyzhou/tinynode/json"
	"sync"
)

/*
 * client face
 * @author <AndyZhou>
 * @mail <diudiu8848@163.com>
 */

//global variable
var (
	_client *Client
	_clientOnce sync.Once
)

//face info
type Client struct {
	c *client.Client
}

//get single instance
func GetClient() *Client {
	_clientOnce.Do(func() {
		_client = NewClient()
	})
	return _client
}

//construct
func NewClient() *Client {
	this := &Client{
		c: client.NewClient(),
	}
	return this
}

//quit
func (f *Client) Quit() {
	f.c.Quit()
}

//api for node
//get nodes
func (f *Client) GetNodes(
		groups ...string,
	) ([]*json.NodeInfo, error) {
	//get monitor address
	monitorsAddr := f.GetMonitorAddr()
	if monitorsAddr == nil || len(monitorsAddr) <= 0 {
		return nil, errors.New("can't get monitor address")
	}

	//loop query
	for _, addr := range monitorsAddr {
		monitor, _ := f.c.GetRpc().GetMonitor(addr)
		if monitor != nil {
			nodesInfo, _ := monitor.GetNodesInfo(groups...)
			if nodesInfo != nil {
				return nodesInfo, nil
			}
		}
	}
	return nil, errors.New("no any node info")
}

//sync node info to monitor
func (f *Client) SyncNode(
		req *json.SyncNodeReq,
		monitorAddresses ...string,
	) error {
	var (
		err error
	)
	//check
	if monitorAddresses == nil {
		monitorAddresses = f.GetMonitorAddr()
	}
	if monitorAddresses == nil {
		return errors.New("no any monitor address")
	}

	//cast to all monitors
	for _, addr := range monitorAddresses {
		monitor, _ := f.GetMonitor(addr)
		if monitor == nil {
			continue
		}
		err = monitor.SyncNodeInfo(req)
	}
	return err
}

//gen new node sync req
func (f *Client) GenSyncNodeReq() *json.SyncNodeReq {
	return json.NewSyncNodeReq()
}

//set cb for node notify
//call cb by stream way
func (f *Client) SetCBForNodeNotify(cb func(info *json.NodeInfo) error) {
	if cb == nil {
		return
	}
	f.c.GetRpc().SetCBForNodeNotify(cb)
}

//api for monitor
//get current monitor addr
func (f *Client) GetMonitorAddr() []string {
	return f.c.GetRpc().GetMonitorAddr()
}

//get monitor
func (f *Client) GetMonitor(address string) (*client.Monitor, error) {
	return f.c.GetRpc().GetMonitor(address)
}

//remove monitor
func (f *Client) RemoveMonitor(address string) error {
	return f.c.GetRpc().RemoveMonitor(address)
}

//add monitor address
func (f *Client) AddMonitor(address ...string) error {
	return f.c.GetRpc().AddMonitor(address...)
}