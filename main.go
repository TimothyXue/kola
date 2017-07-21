package main

import (
	"fmt"
	"kola/client"
	"kola/server"
)

func main() {
	cmd := parseCmd()
	if cmd.versionOption {
		fmt.Printf("The current Kola version is %s. \n", KOLAVERSION)
	} else if cmd.helpOption {
		printUsage()
	} else if cmd.startClient {
		client.StartClient()
	} else if cmd.startServer {
		server.StartServer()
	} else {
		printUsage()
	}
}
