package main

import (
	"cmstopad/common"
	"flag"
	"fmt"
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync"
	"time"
)

var (
	Num  int
	Conn int
	Addr string
)

const RS = 10000

func init() {
	flag.IntVar(&Num, "num", 10, "")
	flag.IntVar(&Conn, "conn", 1, "")
	flag.StringVar(&Addr, "addr", ":3737", "")
	flag.Parse()
}

func main() {

	clients := make([]*rpc.Client, 0, 10)
	//创建并行多个连接
	for i := 0; i < Conn; i++ {
		client, err := jsonrpc.Dial("tcp", Addr)
		if err != nil {
			log.Fatalln(err)
		}
		clients = append(clients, client)
	}

	index := 0

	//等待所有goroutine结束
	var wg sync.WaitGroup

	start := time.Now()

	for i := 0; i < Num; i++ {
		wg.Add(1)

		var c *rpc.Client
		if index < len(clients) {
			c = clients[index]
			index++
		} else {
			index = 0
			c = clients[index]
		}

		go func(cli *rpc.Client) {
			defer wg.Done()

			log.Println("goroutine start...")
			st := time.Now()

			var reply *common.DispResponseData = new(common.DispResponseData)
			//var request *common.DispRequestData = &common.DispRequestData{2, 0, 0, 1, ""}

			for n := 0; n < RS; n++ {
				var adpos_id uint32 = uint32(n%10)
				var request *common.DispRequestData = &common.DispRequestData{adpos_id, 0, 0, 1, ""}
				//err = client.Call("RetrievalRPC.QueryAds", &request, &reply)
				if err := cli.Call("RetrievalRPC.QueryAds", request, reply); err != nil {
					log.Println(err)
				}
			}

			log.Println(time.Now().Sub(st))
		}(c)
	}
	wg.Wait()

	total := RS * Num
	secs := time.Now().Sub(start) / 1000000000

	/*	fmt.Printf("concurrency: %d\n", Num)
		fmt.Printf("total: %d\n", total)
		fmt.Printf("seconds: %d\n", secs)
		fmt.Printf("qps: %d\n", total/int(secs))
		fmt.Printf("num:%d,conn:%d\n", Num,Conn)*/

	fmt.Printf("%d|%d|%d|%d|%d\n", Num, Conn, total, secs, total/int(secs))
}
