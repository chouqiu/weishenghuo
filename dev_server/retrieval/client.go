//  用来测试和  管理服务的客户端

package main

import (
	"flag"
	"fmt"
	"net/rpc"
	"time"
)

func do_client() int {
	addr := *g_disp_addr
	sock := "tcp"

	// client
	client, err := rpc.Dial(sock, addr)
	if err != nil {
		client, err = try_to_connect(sock, addr)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			return 1
		}
	}
	defer client.Close()

	if flag.NArg() > 0 {
		switch flag.Arg(0) {
		case "close":
			cmd_close(client)
		case "status":
			cmd_status(client)
		}
	}
	return 0
}

func try_to_connect(network, address string) (client *rpc.Client, err error) {
	t := 0
	for {
		client, err = rpc.Dial(network, address)
		if err != nil && t < 1000 {
			time.Sleep(10 * time.Millisecond)
			t += 10
			continue
		}
		break
	}

	return
}

//-------------------------------------------------------------------------
// commands
//-------------------------------------------------------------------------

func cmd_status(c *rpc.Client) {
	//fmt.Printf("%s\n", client_status(c, 0))
	fmt.Println("cmd status")
}

func cmd_close(c *rpc.Client) {
	client_close(c, 0)
}
