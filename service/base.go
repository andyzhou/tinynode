package service

import (
	"fmt"
	"github.com/andyzhou/tinynode/cmd"
	"github.com/andyzhou/tinynode/define"
)

//base face
type BaseService struct {
}

//print all usage
func (f *BaseService) PrintAll() {
	f.PrintServiceUsage()
}

//print service usage
func (f *BaseService) PrintServiceUsage() {
	fmt.Println("usage:")
	//for chunk
	fmt.Printf(" %v -%v <rpc port> -%v <log path>\n",
		define.AppName,
		cmd.NameOfPort,
		cmd.NameOfLog,
	)
	fmt.Println()
}