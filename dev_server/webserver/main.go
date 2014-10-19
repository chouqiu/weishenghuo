// Copyright 2014 . All rights reserved.
// create by apple 2014.10.16
//

package main

import (
	"errors"
	"flag"
	//	"github.com/VividCortex/godaemon"
	"github.com/golang/glog"
	"net/http"
	"regexp"
)

//服务器文件根目录
const (
	ServerRootDir = "/var/www/www.targetingx.com/disp/"
)

//当传入参数错误时，返回该错误
//参数错误包括
//1. 参数缺失
//2. 参数不合法
var ErrorInputParameter = errors.New("parameter incorrect")

//访问路径正则匹配
//访问的目录只能是
//1. ping 包括点击统计
//2. disp 内容请求
//3. img/js/html 静态页面请求  这部分 可以考虑放 apache之类的
var validPath = regexp.MustCompile("^/(ping|disp|img|js|html)/([a-zA-Z0-9._/]+)")

//调用相应的处理函数前进行预处理
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//防止程序panic,this is a safe handle.
		defer func() {
			if e, ok := recover().(error); ok {
				http.Error(w, "internal error.", http.StatusInternalServerError)
				glog.Errorf("WARN: panic in %v. - %v", fn, e)
			}
		}()
		//end

		//记录访问日志
		glog.V(0).Infof("%s\t%s\t%s\t%s\t%s\t%s", r.Method, r.RemoteAddr, r.URL.Path, r.Proto, r.Referer(), r.UserAgent())

		//目前只接受get请求，其余请求忽略
		switch r.Method {
		case "GET":
			m := validPath.FindStringSubmatch(r.URL.Path)
			if m == nil {
				http.NotFound(w, r)
				return
			}

			fn(w, r, r.URL.Path)
			glog.Flush()
			return
		default:
			http.Error(w, "method not support...", http.StatusInternalServerError)
			glog.Flush()
			return
		}
	}
}

//g_retrieval_service 为访问的retrieval的ip:port，ip为空时表示本机访问
//g_web_service 程序要监听的http端口，默认为80端口，必须使用root权限才能运行。普通用户只能监听1024以上的端口
//g_ping_service pingserver监听的端口
//g_server_file_dir 网站根目录
//g_aid_cnt 最多每次返回的广告数目，目前为1
var (
	g_retrieval_service = flag.String("retrieval_addr", ":3737", "address for retrieval socket")
	g_web_service       = flag.String("web_service", ":80", "address for http socket")
	g_ping_service      = flag.String("ping_service", ":1234", "address for ping socket")
	g_server_file_dir   = flag.String("server_dir", ServerRootDir, "dir for server")
	g_aid_cnt           = flag.Uint64("aid_cnt", 1, "max aid required")
	g_cache_period      = flag.Uint64("cache_period", 600, "max aid required")
)

func main() {
	//godaemon.MakeDaemon(&godaemon.DaemonAttr{})

	flag.Parse()

	http.HandleFunc("/ping/", makeHandler(MyPingHandler.HandlePing))
	http.HandleFunc("/disp/", makeHandler(MyDispHandler.HandleDisplay))
	//	http.HandleFunc("/img/", makeHandler(StaticDirHandler))
	//	http.HandleFunc("/html/", makeHandler(StaticDirHandler))

	err := http.ListenAndServe(*g_web_service, nil)

	if err != nil {
		panic(err)
	}
}
