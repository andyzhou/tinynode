package main

import (
	"github.com/andyzhou/tinylib/util"
	"github.com/andyzhou/tinynode/cmd"
	"github.com/andyzhou/tinynode/define"
	"github.com/andyzhou/tinynode/service"
	"github.com/urfave/cli"
	"log"
	"os"
)

/*
 * main service face
 * - used for daemon service
 */

// global variable
var (
	logger *util.LogService
)

// start relate service
func start(c *cli.Context) error {
	//init cmd conf
	cmd.InitRunCmdCfg(c)

	//get service key config
	serviceCfg := cmd.GetServiceCfg(c)
	logPath := serviceCfg.LogPath
	if logPath == "" {
		logPath = define.DefaultLogPath
	}

	//register log service
	logger = util.NewLogService(logPath, define.AppName)

	//start core service
	coreDaemon := service.NewCoreDaemon()
	coreDaemon.RunDaemon(c)
	return nil
}

func main() {
	//get command flags
	flags := cmd.Flags()

	//init app
	app := &cli.App{
		Name: define.AppName,
		Action: func(c *cli.Context) error {
			return start(c)
		},
		Flags: flags,
	}

	//start app
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("main run err: %v", err)
		return
	}

	//clean up
	if logger != nil {
		logger.Close()
	}
	log.Printf("%v run succeed!\n", define.AppName)
}
