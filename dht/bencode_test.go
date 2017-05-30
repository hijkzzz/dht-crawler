package dht

import (
	"testing"
)

func test_EncodeBencode(t *testing.T) {
	var m = map[string]interface{}{"1": "a", "2": "b"}
	a, error := EncodeBencode(m)
	if error == nil {
		//fmt.Println(string(a) + "--" + string(error))
	}
}
