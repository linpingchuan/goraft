package main

import (
	"flag"
	"goraft"
	"log"
	"net"
	"strings"

	"google.golang.org/grpc"
)

var (
	peerString string
	nodeName   string
	listenPort string

	peers []string
)

func init() {
	flag.StringVar(&nodeName, "node", "", "the node name has to be unique")
	flag.StringVar(&listenPort, "port", ":7000", "listen port")
	flag.StringVar(&peerString, "peers", "", "peer address, seperated by comma")
}

func main() {
	flag.Parse()
	parseArgs()

	l, err := net.Listen("tcp", listenPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	goraft.RegisterRaftServer(s, &goraft.Raft{})
	s.Serve(l)
}

func parseArgs() {
	if nodeName == "" {
		panic("nodeName cannot be empty")
	}

	if peerString == "" {
		panic("peerString cannot be empty")
	}
	peers = strings.Split(peerString, ",")
}
