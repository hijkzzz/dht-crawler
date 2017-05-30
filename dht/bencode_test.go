package dht

import (
	"fmt"
	"testing"
)

func test_encodeBencode(t *testing.T) {
	var m = map[string]interface{}{"1": "a", "2": "b"}
	a, error := encodeBencode(m)
	if error == nil {
		fmt.Println(string(a) + "--" + fmt.Sprintf("%v", error))
	}
}
