package main

import (
	"github.com/hundred666/BlockchainInfo/spider"
	"fmt"
	"time"
	"github.com/hundred666/BlockchainInfo/block"
)

func main() {
	hashChan := make(chan string, 100)
	outBlockChan := make(chan string, 100)
	failedBlock := make(chan string, 5)
	b := block.Block{
		Hash:        hashChan,
		OutBlock:    outBlockChan,
		FailedBlock: failedBlock,
	}

	failedCount := 0
	flag := false

	b.Hash <- "00000000000000000024c244f9c7d1cc0e593a7a4aa31c1ee2ef35206934bfff"
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case v := <-b.OutBlock:
			go spider.ParseBlock(v, &b)
		case r := <-b.Hash:
			go b.GetBlock(r)
		case v := <-b.FailedBlock:
			failedCount++
			time.Sleep(5 * time.Minute)
			go spider.ParseBlock(v, &b)
		default:

		}

		select {
		case <-ticker.C:
			fmt.Println("timeout")
			flag = true
		default:
		}

		if flag || failedCount == 5 {
			break
		}
	}

	close(b.Hash)
	close(b.OutBlock)

	//b.Stop <- 0

}
