package core

import (
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

func ProcessCommand(conn net.Conn) error {
	lcmd, err := readCommand(conn)
	if err != nil {
		return err
	}
	if len(lcmd) == 0 {
		conn.Write([]byte("no command received"))
	}
	lresp := strings.Fields(lcmd)
	lkvdbCmd := &KVDBCmd{
		Cmd:  lresp[0],
		Args: lresp[1:],
	}
	lkvdbCmd.Parse(conn)
	return nil
}

func readCommand(lconn net.Conn) (string, error) {
	var lbuf []byte = make([]byte, 512)
	n, err := lconn.Read(lbuf)
	if err != nil {
		return "", err
	}
	return string(lbuf[:n]), nil
}

func (kvcmd *KVDBCmd) Parse(w io.ReadWriter) {
	switch kvcmd.Cmd {
	case "PING":
		w.Write([]byte("Hello, your ping is a success.."))
	case "SET":
		if len(kvcmd.Args) > 3 {
			w.Write([]byte("incorrect syntax"))
			return
		}
		handleSet(kvcmd.Args, w)
	case "GET":
		if len(kvcmd.Args) == 0 || len(kvcmd.Args) > 1 {
			w.Write([]byte("incorrect syntax"))
			return
		}
		handleGet(kvcmd.Args, w)
	default:
		w.Write([]byte("unsupported command"))
		return
	}
}

func handleSet(args []string, w io.ReadWriter) {
	var ldurInSec int64 = -1
	var err error
	if len(args) == 3 {
		ldurInSec, err = strconv.ParseInt(args[2], 10, 64)
		if err != nil {
			w.Write([]byte("duration provided for the command is out of range"))
			return
		}
		ldurInSec = ldurInSec * 1000
	} else {
		ldurInSec = -1
	}
	lkvobj := NewKVObj(args[1], ldurInSec)
	Put(args[0], lkvobj)
	w.Write([]byte("key is set"))
}

func handleGet(args []string, w io.ReadWriter) {
	lresp := Get(args[0])
	if lresp == nil {
		w.Write([]byte("key does not exist"))
		return
	}
	if lresp.ExpiresAt > 0 {
		if lresp.ExpiresAt < time.Now().UnixMilli() {
			w.Write([]byte("value has expired"))
			return
		}
	}
	w.Write([]byte(lresp.Value))
}
