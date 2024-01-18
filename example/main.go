package main

import (
	"github.com/andyzhou/tinynode"
	"github.com/andyzhou/tinynode/define"
	"github.com/andyzhou/tinynode/json"
	"log"
	"sync"
)

/*
 * example code
 */

const (
	MonitorAddr = "localhost:5300"
	NodeTag = "test"
	NodeAddr = "localhost:7100"
	NodeGroup = "chat"
)

//cb for node notify
func cbForNodeNotify(info *json.NodeInfo) error {
	log.Printf("info:%v\n", info)
	return nil
}

func main() {
	var (
		wg sync.WaitGroup
	)

	//init client
	c := tinynode.NewClient()

	//set cb for node notify
	c.SetCBForNodeNotify(cbForNodeNotify)

	//add monitor
	c.AddMonitor(MonitorAddr)

	//add wg
	wg.Add(1)
	log.Printf("example running...\n")

	//sync node info
	req := c.GenSyncNodeReq()
	req.Tag = NodeTag
	req.Address = NodeAddr
	req.Group = NodeGroup
	req.Stat = define.NodeStatOfBusy

	err := c.SyncNode(req)
	if err != nil {
		log.Printf("sync node failed, err:%v\n", err.Error())
	}

	//wait
	wg.Wait()
}