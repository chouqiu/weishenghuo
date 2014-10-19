// Copyright 2014 . All rights reserved.
// create by apple 2014.10.16
//

package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"net/http"
	"net/rpc/jsonrpc"
	//	"regexp"
	"strconv"
	"text/template"
	"weishenghuo/common"
)

const (
	InternalConnectionQueryError  = "internal connection err for query ads"
	InternalConnectionReportError = "internal connection err for disp report"
	InternalConnectionReturnError = "internal connection err for return"
	MaxInputCallbackLength        = 20 //callback长度最大不能超过20个字符
	DefaultCallbackName           = "callback"
)

var MyDispHandler DispHandler = DispHandler{}

type DispHandler struct {
}

//曝光请求
func (disp_handler DispHandler) HandleDisplay(w http.ResponseWriter, r *http.Request, target_page string) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}

	request_cgi := m[2]

	switch request_cgi {
	case "req.fcg":
		disp_handler.handle_display(w, r)
	default:
		http.NotFound(w, r)
		return
	}
}

//获取曝光请求参数并访问retrieval获取新闻内容
//广告会带上traceid
func (disp_handler DispHandler) handle_display(w http.ResponseWriter, r *http.Request) {

	//解析参数，需要广告位和callback
	locate_id, err := disp_handler.parse_input(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//请求曝光
	client, err := jsonrpc.Dial("tcp", *g_retrieval_service)

	if err != nil {
		if client, err = jsonrpc.Dial("tcp", *g_retrieval_service); err != nil {
			glog.Error("call jsonrpc.Dial error")
			http.Error(w, InternalConnectionQueryError, http.StatusInternalServerError)
			return
		}
	}

	defer client.Close()

	//rpc请求新闻
	var request common.NewsListRequest = common.NewsListRequest{}
	var reply common.NewsListResponse = common.NewsListResponse{}

	request.LocationId = locate_id

	err = client.Call("RetrievalRPC.QueryNews", &request, &reply)

	if err != nil {
		glog.Error("call RetrievalRPC.QueryNews error")
		http.Error(w, InternalConnectionQueryError, http.StatusInternalServerError)
		return
	}

	//reply中带上traceid
	result, err := json.Marshal(reply)
	if err != nil {
		http.Error(w, InternalConnectionReturnError, http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(result))
	return
}

func (disp_handler DispHandler) parse_input(r *http.Request) (locate_id uint32, err error) {

	//先解析获取数据
	err = r.ParseForm()

	if err != nil {
		return
	}

	locate_id = 0
	var tmp uint64 = 0

	str_locate_id, ok := r.Form["locate_id"]

	if ok == false {
		err = ErrorInputParameter
		return
	} else {
		tmp, err = strconv.ParseUint(template.HTMLEscapeString(str_locate_id[0]), 10, 32)
		if err != nil || tmp == 0 {
			err = ErrorInputParameter
			return
		} else {
			locate_id = uint32(tmp)
		}
	}

	err = nil
	return
}
