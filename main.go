package main

import "dht-crawler/dht"

func main() {
	dht := dht.NewDHT("0.0.0.0", 34568)
	dht.Run()
}
