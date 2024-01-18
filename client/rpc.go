package client

import (
	"errors"
	"github.com/andyzhou/tinynode/json"
	"sync"
)

/*
 * rpc client face
 * @author <AndyZhou>
 * @mail <diudiu8848@163.com>
 * - send request to service
 * - real time receive node sync pass stream way
 * - support local and multi modes
 */

//face info
type Rpc struct {
	monitorMap map[string]*Monitor //addr -> *Monitor
	cbForNodeNotify func(info *json.NodeInfo) error
	sync.RWMutex
}

//construct
func NewRpc() *Rpc {
	this := &Rpc{
		monitorMap: map[string]*Monitor{},
	}
	return this
}

//quit
func (f *Rpc) Quit() {
	f.Lock()
	defer f.Unlock()
	for _, v := range f.monitorMap {
		v.Quit()
	}
	f.monitorMap = map[string]*Monitor{}
}

//get monitor address
func (f *Rpc) GetMonitorAddr() []string {
	result := make([]string, 0)
	for addr, _ := range f.monitorMap {
		if addr != "" {
			result = append(result, addr)
		}
	}
	return result
}

//remove monitor by addr
func (f *Rpc) RemoveMonitor(addr string) error {
	//check
	if addr == "" {
		return errors.New("invalid parameter")
	}
	//get monitor
	monitor, _ := f.GetMonitor(addr)
	if monitor != nil {
		monitor.Quit()
		//remove with locker
		f.Lock()
		defer f.Unlock()
		delete(f.monitorMap, addr)
	}
	return nil
}

//get monitor by addr
func (f *Rpc) GetMonitor(addr string) (*Monitor, error) {
	//check
	if addr == "" {
		return nil, errors.New("invalid parameter")
	}
	f.Lock()
	defer f.Unlock()
	v, ok := f.monitorMap[addr]
	if ok && v != nil {
		return v, nil
	}
	return nil, errors.New("no such monitor")
}

//add monitor server node
func (f *Rpc) AddMonitor(addrSlice ...string) error {
	//check
	if addrSlice == nil || len(addrSlice) <= 0 {
		return errors.New("invalid parameter")
	}
	for _, addr := range addrSlice {
		if addr == "" {
			continue
		}
		//check monitor
		v, _ := f.GetMonitor(addr)
		if v != nil {
			continue
		}

		//init new monitor
		monitor := NewMonitor(addr)
		monitor.SetCBForNodeNotify(f.cbForNodeNotify)

		//sync into map
		f.Lock()
		f.monitorMap[addr] = monitor
		f.Unlock()
	}
	return nil
}

//set cb for node notify
func (f *Rpc) SetCBForNodeNotify(cb func(info *json.NodeInfo) error) {
	if cb == nil {
		return
	}
	f.cbForNodeNotify = cb
}
