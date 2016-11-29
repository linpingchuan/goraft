package goraft

import (
	"log"
	"net"
	"testing"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func TestRPC(t *testing.T) {
	listenPort := ":7001"
	go startServer(listenPort)

	// make client
	cli := getClient(listenPort)

	cxt := context.Background()
	aeReq := &AERequest{LeaderID: "node1", Term: 12}
	aeResp, err := cli.AppendEntries(cxt, aeReq)
	if err != nil {
		panic(err)
	}
	log.Printf("%#v", aeResp)
}

func startServer(port string) {
	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	RegisterRaftServer(s, &Raft{})
	s.Serve(l)
}

func getClient(port string) RaftClient {
	conn, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		log.Printf(`connect to service failed: %v`, err)

	}
	cli := NewRaftClient(conn)
	return cli
}
