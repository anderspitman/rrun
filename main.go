package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"os/exec"
)

func main() {
	netAddr := flag.String("addr", "localhost", "Net address")
	port := flag.Int("port", 9001, "Port")
	flag.Parse()

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *netAddr, *port))
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	cmdSizeBuf := make([]byte, 4)
	_, err := conn.Read(cmdSizeBuf)
	if err != nil {
		return
	}
	cmdSize := binary.BigEndian.Uint32(cmdSizeBuf)

	cmdBuf := make([]byte, cmdSize)
	_, err = conn.Read(cmdBuf)
	if err != nil {
		return
	}

	cmdStr := string(cmdBuf)
	logLine := fmt.Sprintf("\t%s\t%s", conn.RemoteAddr(), cmdStr)
	log.Printf(logLine)

	cmd := exec.Command("/bin/sh", "-c", string(cmdBuf))
	cmd.Stdin = conn
	cmd.Stdout = conn
	cmd.Run()
}
