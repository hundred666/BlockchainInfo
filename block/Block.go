package block

import (
	"io/ioutil"
	"net/http"
	"fmt"
	"time"
)

const RawBlock = "https://blockchain.info/rawblock/"

type Block struct {
	Hash        chan string
	OutBlock    chan string
	CurrBlock   string
	FailedBlock chan string
}

func (b *Block) GetBlock(hash string) {

	url := RawBlock + hash
	fmt.Println("fetching block ", hash, " at time ", time.Now())
	resp, err := http.Get(url)
	if err != nil {
		b.FailedBlock <- hash
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		b.FailedBlock <- hash
		return
	}
	b.OutBlock <- string(body)
	resp.Body.Close()

}
