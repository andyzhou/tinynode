package main

import (
	"github.com/andyzhou/tinynode"
	"log"
	"sync"
)

/*
 * example code
 */

const (
	MonitorAddr = "localhost:5300"
)

func main() {
	var (
		wg sync.WaitGroup
	)

	//init client
	c := tinynode.NewClient()

	//add monitor
	c.AddMonitor(MonitorAddr)

	//add wg
	wg.Add(1)
	log.Printf("example running...\n")

	//wait
	wg.Wait()
}