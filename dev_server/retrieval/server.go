package main

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
)

func do_server() int {

	disp_addr := *g_disp_addr

	g_daemon = new_daemon(disp_addr)

	rpc.Register(g_retrieval_rpc)

	//protobuf rpc自身封装好了，无法修改
	//go common.ListenAndServeCmstopAdService("tcp", *g_index_addr, g_retrieval_rpc)

	g_daemon.loop()
	return 0
}

//-------------------------------------------------------------------------
// daemon
//-------------------------------------------------------------------------

type daemon struct {
	listener net.Listener //
	cmd_in   chan int     // 这个通道用于 管理 server
}

func new_daemon(address string) *daemon {
	var err error

	d := new(daemon)
	d.listener, err = net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}

	d.cmd_in = make(chan int, 1)

	return d
}

const (
	daemon_close = iota
)

func (this *daemon) loop() {
	conn_in := make(chan net.Conn)

	// 启动 更新 goroutine
	//go CheckIndexUpdate()
	go func() {
		for {
			time.Sleep(5 * time.Second) // 等待5s检查一次全量
			g_check_total_expire <- true
		}
	}()

	// 接受cgi  链接
	go func() {
		for {
			c, err := this.listener.Accept()
			if err != nil {
				continue
			}
			conn_in <- c
		}
	}()

	for {
		// handle connections or server CMDs (currently one CMD)
		select {
		case c := <-conn_in:
			go jsonrpc.ServeConn(c)
		//	runtime.GC()
		case cmd := <-this.cmd_in:
			switch cmd {
			case daemon_close:
				return
			}
		case <-g_check_total_expire:
			CheckNewsUpdate()
			//	case update_request := <-g_index_incr_chan:
			//		ServerUpdateIndex(&update_request)
		}
	}
}

func (this *daemon) close() {
	this.cmd_in <- daemon_close
}

var g_daemon *daemon

//-------------------------------------------------------------------------
// server_* functions for rpc
//
//-------------------------------------------------------------------------

func server_close(notused int) int {
	g_daemon.close()
	return 0
}
