package main

import (
	"fmt"
	"time"
)

func worker(done chan struct{}) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("TIK")
		case <-done:
			fmt.Println("stopping ...")
			return
		}
	}

}

func main() {
	done := make(chan struct{})

	go worker(done)
	time.Sleep(5 * time.Second)
	close(done)

	time.Sleep(2 * time.Second)
}
