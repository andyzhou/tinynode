package service

import (
	"fmt"
	"github.com/andyzhou/tinylib/util"
	"github.com/andyzhou/tinynode/cmd"
	"github.com/andyzhou/tinynode/define"
	"github.com/andyzhou/tinynode/service/face"
	"github.com/urfave/cli"
	"sync"
)

/*
 * core daemon service
 * @author <AndyZhou>
 * @mail <diudiu8848@163.com>
 */

//service info
type CoreDaemon struct {
	signal *util.Signal
	shutDownChan chan bool
	wg sync.WaitGroup
	BaseService
}

//construct
func NewCoreDaemon() *CoreDaemon {
	this := &CoreDaemon{
		signal: util.NewSignal(),
		shutDownChan: make(chan bool, 1),
	}
	this.interInit()
	return this
}

//run core daemon server
func (d *CoreDaemon) RunDaemon(c *cli.Context) {
	//get sub service config from command flag
	serviceCfg := cmd.GetServiceCfg(c)

	//check args config
	if serviceCfg == nil ||
		serviceCfg.RpcPort <= 0 {
		d.PrintServiceUsage()
		return
	}

	//wait group add
	d.wg.Add(1)

	//init inter data
	//run in son process
	d.initInterData(serviceCfg)

	//init and start rpc service
	rpcService := NewRpcService(serviceCfg)

	//print info and wait
	fmt.Printf("start %v on port %v..\n", define.AppName, serviceCfg.RpcPort)
	d.wg.Wait()

	//clean up
	rpcService.Quit()
	fmt.Printf("stop %v on port `%v` ..\n", define.AppName, serviceCfg.RpcPort)
}

//init inter data
func (d *CoreDaemon) initInterData(cfg *cmd.ServiceCfg) {
	//check
	if cfg == nil {
		return
	}
	//init inter face
	face.GetInterFace()
}

//cb for daemon signal shutdown
func (d *CoreDaemon) cbForDaemonShutDown() {
	//inter data quit
	face.GetInterFace().Quit()

	//wait group done
	d.wg.Done()
}

//inter init
func (d *CoreDaemon) interInit() {
	//register signal quit
	d.signal.RegisterShutDownChan(d.shutDownChan, d.cbForDaemonShutDown)

	//monitor signal
	d.signal.MonSignal()
}