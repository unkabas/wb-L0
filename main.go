package main

import (
	"flag"
	"fmt"
	"l0/cmd"
	"l0/iternal/config"
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
