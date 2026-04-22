package main

import (
	"go-load-balancer/proxies"
	"go-load-balancer/servers"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go servers.Http(&wg)
	go proxies.Http(&wg)

	wg.Wait()

}
