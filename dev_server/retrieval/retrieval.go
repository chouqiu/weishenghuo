// Copyright 2014 . All rights reserved.
// create by apple 2014.10.15
//

package main

import (
	"flag"
	"fmt"
	"weishenghuo/common"
	//	"github.com/VividCortex/godaemon"
	"os"
)

//retrieval 新闻索引 结构体
type NewsIndexProxy struct {
	news_index map[uint32]common.NewsCluster
}

type RetrievalRPC struct {
	news_index_proxy *NewsIndexProxy //

}

func (this *RetrievalRPC) QueryNews(disp_request *common.NewsListRequest, disp_response *common.NewsListResponse) error {
	*disp_response = ServerQueryNews(&this.news_index_proxy.news_index, disp_request)
	return nil
}

//接受indexer的增量更新信息
//这个函数是proto里面定义的，rpc必须实现一个一样的函数才能用于protorpc
//func (this *RetrievalRPC) Update(update_request *common.CmstopAdRequest, update_response *common.CmstopAdResponse) error {
//	g_index_incr_chan <- (*update_request)
//	update_response.RetCode = proto.Uint32(0)
//	return nil
//}

//g_disp_addr 作为曝光服务端时监听的端口
//g_index_file_dir 全量索引文件的目录
//g_is_server 绝对采用客户端还是服务端的模式启动程序
//g_index_addr 作为索引服务端接受增量信息时的监听端口
var (
	g_disp_addr        = flag.String("disp_addr", ":3737", "address for tcp socket")
	g_index_file_dir   = flag.String("index", "../index/", "dir for index update")
	g_is_server        = flag.Bool("s", true, "run a server instead of a client")
	g_index_addr       = flag.String("index_addr", ":3738", "address for tcp socket")
	g_click_url_prefix = flag.String("click_url_prefix", "/", "click url prefix")
)

var g_retrieval_rpc *RetrievalRPC
var g_check_total_expire chan bool

func main() {
	//godaemon.MakeDaemon(&godaemon.DaemonAttr{})
	fmt.Println("starting retrieval")

	flag.Parse()

	//------------------------------------------
	//retrieval rpc proxy 初始化
	//------------------------------------------
	g_retrieval_rpc = new(RetrievalRPC)
	g_retrieval_rpc.news_index_proxy = new(NewsIndexProxy)
	g_retrieval_rpc.news_index_proxy.news_index = make(map[uint32]common.NewsCluster)

	//----------------------------------------------
	//chan初始化
	//----------------------------------------------
	g_check_total_expire = make(chan bool)

	var retval int
	if *g_is_server {
		retval = do_server()
	} else {
		retval = do_client()
	}
	os.Exit(retval)

}

//作为服务端，返回news信息
func ServerQueryNews(news_index *map[uint32]common.NewsCluster, dis_request *common.NewsListRequest) (dis_response common.NewsListResponse) {

	//请求广告位必须有效且有新闻数据
	var location_id uint32 = dis_request.LocationId

	//读取列表
	_, ok := (*news_index)[location_id]

	if ok == false {
		dis_response = common.NewsListResponse{
			Ret:  1, //不存在
			Data: common.NewsList{},
		}
		return dis_response
	}

	dis_response = common.NewsListResponse{
		Ret:  0,
		Data: common.NewsList{},
	}

	return dis_response
}
