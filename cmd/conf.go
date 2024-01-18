package cmd

import "github.com/urfave/cli"

//global variable
var (
	RunCmdConf ServiceCfg
)

//init run command cfg
func InitRunCmdCfg(c *cli.Context) {
	//get relate cmd conf
	serviceCfg := GetServiceCfg(c)

	//setup global variable
	RunCmdConf = *serviceCfg
}

//command flags
func Flags() []cli.Flag  {
	return []cli.Flag{
		&cli.IntFlag{Name: NameOfPort, Usage: "rpc service port"},
		&cli.StringFlag{Name: NameOfLog, Usage: "log path"},
	}
}

//get service config
func GetServiceCfg(c *cli.Context) *ServiceCfg {
	cfg := &ServiceCfg{
		RpcPort: c.Int(NameOfPort),
		LogPath: c.String(NameOfLog),
	}
	return cfg
}