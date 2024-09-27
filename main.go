package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"strings"

	"github.com/arunraghunath/kvdb/config"
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
		lcmd, err := processCommand(lconn)
		if err != nil {
			println(err)
		} else {
			fmt.Println(lcmd)
			lconn.Write([]byte(lcmd[0]))
		}
	}
}

func readCommand(lconn net.Conn) (string, error) {
	var lbuf []byte = make([]byte, 512)
	n, err := lconn.Read(lbuf)
	if err != nil {
		return "", err
	}
	println("printing readcommand" + string(lbuf[:n]))
	return string(lbuf[:n]), nil
}

func processCommand(conn net.Conn) ([]string, error) {
	lcmd, err := readCommand(conn)
	if err != nil {
		return nil, err
	}
	if len(lcmd) == 0 {
		return nil, errors.New("no command received")
	}
	lresp := strings.Fields(lcmd)
	return lresp, nil
}
