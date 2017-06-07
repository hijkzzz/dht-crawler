package dht

import (
	"testing"
)

func Test_run(t *testing.T) {
	var seed = "@hujian:@liujianbiao:@wangpeijia"
	dht := NewDHT("127.0.0.1", 34567, seed)
	dht.Run()
}
