package main

import (
	"net/rpc"
)

// wrapper for: server_close

type Args_close struct {
	Arg0 int
}
type Reply_close struct {
	Arg0 int
}

func (this *RetrievalRPC) RPC_close(args *Args_close, reply *Reply_close) error {
	reply.Arg0 = server_close(args.Arg0)
	return nil
}
func client_close(cli *rpc.Client, Arg0 int) int {
	var args Args_close
	var reply Reply_close
	args.Arg0 = Arg0
	err := cli.Call("RetrievalRPC.RPC_close", &args, &reply)
	if err != nil {
		panic(err)
	}
	return reply.Arg0
}
