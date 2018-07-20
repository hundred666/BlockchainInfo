package transaction

import (
	"fmt"
	"strings"
)

type Transaction struct {
	Index       float64
	Hash        string
	Size        float64
	InputValue  float64
	OutputValue float64
	Value       float64
	Weight      float64
	SubmitTime  float64
	LockTime    float64
	PackTime    float64
	Height      float64
	ConfirmTime float64
}

func (t *Transaction) ToCSVFormat() ([]string) {
	str := fmt.Sprintf("%.f,%v,%.f,%.f,%.f,%.f,%.f,%.f,%.f,%.f,%.f,%.f", t.Index, t.Hash, t.Size, t.InputValue, t.OutputValue, t.Value, t.Weight, t.SubmitTime, t.LockTime, t.PackTime, t.Height, t.ConfirmTime)
	return strings.Split(str, ",")
}
