package main

import (
	"fmt"
	"forward_proxy"
)

func main() {
	fmt.Println("main init")

	fmt.Println("start forward-proxy")
	go forward_proxy.ForwardPxy()

	select {}
}
