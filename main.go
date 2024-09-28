package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/arunraghunath/kvdb/config"
	"github.com/arunraghunath/kvdb/core"
)

func main() {
	flag.StringVar(&config.Host, "host", "localhost", "Default Host for the TCP Listener")
	flag.StringVar(&config.Port, "port", "6560", "Default Port for the TCP Listener")
	flag.Parse()
	laddr := config.Host + ":" + config.Port

	listener, err := net.Listen("tcp", laddr)
	if err != nil {
		panic(err)
	}
	println("TCP Listener started at localhost on port " + laddr + ". Awaiting client connections..")
	for {
		lconn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		err = core.ProcessCommand(lconn)
		if err != nil {
			fmt.Println(err)
		}
	}
}
