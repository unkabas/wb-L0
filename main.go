package main

import (
	"flag"
	"fmt"
	"github.com/unkabas/wb-L0/cmd"
	"github.com/unkabas/wb-L0/iternal/config"
)

var migrate = flag.Bool("m", false, "migration")

func init() {
	config.LoadEnvs()
	config.ConnectDB()
}
func main() {
	flag.Parse()
	fmt.Println("go")
	if *migrate {
		cmd.Migration()
	}
}
