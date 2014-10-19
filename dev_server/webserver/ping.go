// Copyright 2014 . All rights reserved.
// create by apple 2014.10.16
//

package main

import (
	"net/http"
)

var MyPingHandler PingHandler = PingHandler{}

type PingHandler struct {
}

//曝光和点击  回报
func (ping_handler PingHandler) HandlePing(w http.ResponseWriter, r *http.Request, title string) {

	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}

}
