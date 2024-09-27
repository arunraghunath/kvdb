package core

import (
	"errors"
	"net"
	"strconv"
	"strings"
	"time"
)

func ProcessCommand(conn net.Conn) (string, error) {
	lcmd, err := readCommand(conn)
	if err != nil {
		return "", err
	}
	if len(lcmd) == 0 {
		return "", errors.New("no command received")
	}
	lresp := strings.Fields(lcmd)
	lkvdbCmd := &KVDBCmd{
		Cmd:  lresp[0],
		Args: lresp[1:],
	}
	return lkvdbCmd.Parse()
}

func readCommand(lconn net.Conn) (string, error) {
	var lbuf []byte = make([]byte, 512)
	n, err := lconn.Read(lbuf)
	if err != nil {
		return "", err
	}
	return string(lbuf[:n]), nil
}

func (kvcmd *KVDBCmd) Parse() (string, error) {
	switch kvcmd.Cmd {
	case "PING":
		return "Hello, your ping is a success..", nil
	case "SET":
		if len(kvcmd.Args) > 3 {
			return "", errors.New("incorrect syntax")
		}
		return "", handleSet(kvcmd.Args)
	case "GET":
		if len(kvcmd.Args) == 0 || len(kvcmd.Args) > 1 {
			return "", errors.New("incorrect syntax")
		}
		return handleGet(kvcmd.Args)
	default:
		return "", errors.New("unsupported command")
	}
}

func handleSet(args []string) error {
	var ldurInSec int64 = -1
	var err error
	if len(args) == 3 {
		ldurInSec, err = strconv.ParseInt(args[2], 10, 64)
		if err != nil {
			return errors.New("duration provided for the command is out of range")
		}
		ldurInSec = ldurInSec * 1000
	} else {
		ldurInSec = -1
	}
	lkvobj := NewKVObj(args[1], ldurInSec)
	Put(args[0], lkvobj)
	return errors.New("key is set")
}

func handleGet(args []string) (string, error) {
	lresp := Get(args[0])
	if lresp == nil {
		return "", errors.New("key does not exist")
	}
	if lresp.ExpiresAt > 0 {
		if lresp.ExpiresAt < time.Now().UnixMilli() {
			return "", errors.New("value has expired")
		}
	}
	return lresp.Value, nil
}
