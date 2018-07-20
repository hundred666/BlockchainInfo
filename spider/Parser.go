package spider

import (
	"github.com/hundred666/BlockchainInfo/transaction"
	"encoding/json"
	"encoding/csv"
	"os"
	"fmt"
	"time"
	"github.com/hundred666/BlockchainInfo/block"
)

func ParseBlock(block string, b *block.Block) {

	var result map[string]interface{}
	err := json.Unmarshal([]byte(block), &result)
	if err != nil {
		return
	}
	currHash := result["hash"].(string)
	b.CurrBlock = currHash
	fmt.Println("parsing block ", result["hash"], " at time ", time.Now())
	b.Hash <- result["prev_block"].(string)
	packTime := result["time"].(float64)
	height := result["height"].(float64)
	txRaws := result["tx"].([]interface{})
	txs := make([]transaction.Transaction, 0)
	for _, txRaw := range txRaws {
		txMap := txRaw.(map[string]interface{})
		index := txMap["tx_index"].(float64)
		size := txMap["size"].(float64)
		hash := txMap["hash"].(string)
		submitTime := txMap["time"].(float64)
		weight := txMap["weight"].(float64)
		lockTime := txMap["lock_time"].(float64)
		confirmTime := packTime - submitTime

		inputs := txMap["inputs"].([]interface{})
		outs := txMap["out"].([]interface{})
		value, inputValue, outValue := parseFee(inputs, outs)
		if value == -1 {
			continue //coinbase tx;continue
		}

		tx := transaction.Transaction{
			Index:       index,
			Hash:        hash,
			Size:        size,
			InputValue:  inputValue,
			OutputValue: outValue,
			Value:       value,
			Weight:      weight,
			SubmitTime:  submitTime,
			LockTime:    lockTime,
			PackTime:    packTime,
			Height:      height,
			ConfirmTime: confirmTime,
		}
		txs = append(txs, tx)
	}
	go SaveFile(txs)

	return
}

func parseFee(inputs []interface{}, outs []interface{}) (fee float64, inputValue float64, outValue float64) {
	inputValue = parseInputValue(inputs)
	if inputValue == -1 {
		return -1, 0, 0
	}
	outValue = parseOutputValue(outs)

	return inputValue - outValue, inputValue, outValue
}

func parseInputValue(inputs []interface{}) (float64) {
	v := 0.0
	for _, inputRaw := range inputs {
		input := inputRaw.(map[string]interface{})
		p, e := input["prev_out"]
		if !e {
			return -1 //coinbase tx
		}
		prevOut := p.(map[string]interface{})
		value := prevOut["value"].(float64)
		v += value
	}
	return v
}

func parseOutputValue(outs []interface{}) (float64) {
	v := 0.0
	for _, outRaw := range outs {
		out := outRaw.(map[string]interface{})
		value := out["value"].(float64)
		v += value
	}
	return v
}

func SaveFile(txs []transaction.Transaction) {
	m.Lock()
	f, err := os.OpenFile("20180618.csv", os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		f, _ = os.Create("20180618.csv")
	}

	w := csv.NewWriter(f) //创建一个新的写入文件流
	txCSV := make([][]string, 0)
	for _, tx := range txs {
		txCSV = append(txCSV, tx.ToCSVFormat())
	}

	w.WriteAll(txCSV) //写入数据
	w.Flush()
	f.Close()
	m.Unlock()
}
