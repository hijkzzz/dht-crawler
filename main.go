package main

import "dht-crawler/dht"

func main() {
	var seed = "@hujian:@liujianbiao:@wangpeijia"
	dht := dht.NewDHT("127.0.0.1", 34567, seed)
	dht.Run()
}
