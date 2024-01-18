package cmd

//flag field name
const (
	NameOfPort     = "port"     //rpc port
	NameOfLog      = "log"      //log path
)

//command config define
type (
	ServiceCfg struct {
		RpcPort  int
		LogPath  string
	}
)
