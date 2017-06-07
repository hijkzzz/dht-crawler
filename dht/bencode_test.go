package dht

import (
	"fmt"
	"strconv"
	"testing"
)

var maps []map[string]interface{}

func Test_DecodeBencode(t *testing.T) {
	fmt.Println("DecodeBencode解码---------->")
	var as []string
	// as = append(as, "d1:ad2:id20:abcdefghij0123456899:info_hash20:mnopqrstuvwxyz1234564:porti6881e5:token8:aoeusnthe1:q13:announce_peer1:t2:aa1:y1:qe")
	// as = append(as, "i3e")
	// as = append(as, "d1:rd2:id20:mnopqrstuvwxyz123456e1:t2:aa1:y1:re")
	// as = append(as, "d1:rd2:idt2:aa1:y1:r")
	as = append(as, "d1:t2:�H1:y1:q1:q9:find_node1:ad2:id40:fc7ae309aad10ad10d1e6770290b8822782b93d86:target40:fc7ae309aad10ad10d1e6770290b8822782b93d8ee")
	for i, a := range as {
		fmt.Println("-----------------------------------------------------------------------------------------------")
		fmt.Println("第 " + strconv.Itoa(i) + " 组数据测试：")
		values, error := decodeBencode([]byte(a))
		maps = append(maps, values)
		if error != nil {
			fmt.Println(error)
		}
		for key := range values {
			fmt.Println(key, values[key])
		}
		fmt.Println("-----------------------------------------------------------------------------------------------")
	}

}

func Test_EncodeBencode(t *testing.T) {
	var m = map[string]interface{}{
		"t": "�h",
		"y": "q",
		"q": "find_node",
		"a": map[string]interface{}{
			"id":     "fc7ae309aad10ad10d1e6770290b8822782b93d8",
			"target": "fc7ae309aad10ad10d1e6770290b8822782b93d8",
		},
	}
	maps = append(maps, m)
	fmt.Println("EecodeBencode编码---------->")
	for i, m := range maps {
		fmt.Println("-----------------------------------------------------------------------------------------------")
		fmt.Println("第 " + strconv.Itoa(i) + " 组数据测试：")
		values, error := encodeBencode(m)

		if error != nil {
			fmt.Println(error)
		}
		fmt.Println(m)
		fmt.Println(string(values))

		fmt.Println("-----------------------------------------------------------------------------------------------")
	}
}
