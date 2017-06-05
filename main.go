package main

import "dht-crawler/dht"

func main() {
	var seed = "@hujian:@liujianbiao:@wangpeijia"
	dht := dht.NewDHT("0.0.0.0", 34568, seed)
	dht.Run()
}
