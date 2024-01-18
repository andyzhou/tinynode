package cmd

//flag field name
const (
	NameOfPort     = "port"     //rpc port
	NameOfMonitors = "monitors" //all addr split by ','
	NameOfLog      = "log"      //log path
)

//command config define
type (
	ServiceCfg struct {
		RpcPort  int
		Monitors string //addr split by ','
		LogPath  string
	}
)
