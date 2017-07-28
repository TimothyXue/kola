package main

import (
	"flag"
	"fmt"
	"os"
)

type KolaCmd struct {
	registerAgent bool
	helpOption    bool
	versionOption bool
	startClient   bool
	startServer   bool
	address       KolaHost

	//args          []string
}

// KolaHost contains the address information to start server
type KolaHost struct {
	host string
	port int
}

func parseCmd() *KolaCmd {
	cmd := &KolaCmd{}
	flag.Usage = printUsage
	flag.BoolVar(&cmd.helpOption, "help", false, "print the guidance")
	flag.BoolVar(&cmd.helpOption, "?", false, "print the guidance")
	flag.BoolVar(&cmd.versionOption, "version", false, "demonstrate the version of Kola")
	flag.BoolVar(&cmd.versionOption, "v", false, "demonstrate the version of Kola")
	flag.BoolVar(&cmd.startClient, "client", false, "start kola client")
	flag.BoolVar(&cmd.startServer, "server", false, "start kola server")
	flag.StringVar(&cmd.address.host, "localhost", "define the host name for Kola service")
	flag.IntVar(&cmd.address.port, "port", 5051, "define the port to connect or start")

	flag.Parse()
	//args := flag.Args()
	return cmd
}

func printUsage() {
	fmt.Printf("Usage: %s [-register] <host> \n", os.Args[0])
}
