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
	//args          []string
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
	flag.Parse()
	//args := flag.Args()
	return cmd
}

func printUsage() {
	fmt.Printf("Usage: %s [-register] <host> \n", os.Args[0])
}
