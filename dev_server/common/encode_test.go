package common

import (
    "testing"
	"encoding/hex"
	"fmt"
)

func Test_Encoding(t *testing.T) {
	source_url := "http://www.qq.com?adpos_id=2&order_id=4&creative_id=1&traceid=29804444444"
	
	encode_bin,err := UrlEncode([]byte(source_url))
	if err != nil {
		fmt.Println("get encode_bin error")
		return
	}
	
	encode_result := hex.EncodeToString(encode_bin)
	
	fmt.Println(encode_result)
	
	encode_hex,err := hex.DecodeString(encode_result)
	
	if err != nil {
		fmt.Println("get encode_hex err")
		return
	}
	
	result, err := UrlDecode(encode_hex)
	if err != nil {
		fmt.Println("get result err")
		return
	}
	
	fmt.Println(string(result))
}
