package main

import (
	"flag"
	"github.com/unkabas/wb-L0/cmd/migration"
	"github.com/unkabas/wb-L0/iternal/config"
)

var migrate = flag.Bool("m", false, "migration")

func init() {
	config.LoadEnvs()
	config.ConnectDB()
}
func main() {
	flag.Parse()
	if *migrate {
		migration.Migration()
	}
}
