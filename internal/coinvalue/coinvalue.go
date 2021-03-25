package coinvalue

import (
	"math/big"

	"github.com/leekchan/accounting"
)

type CoinValue struct {
	Value    string
	Emojicon string
}

func New(value float64) CoinValue {
	cv := CoinValue{}
	ac := accounting.Accounting{Symbol: "", Precision: 2}

	cv.Emojicon = "✅"
	if value < 0 {
		cv.Emojicon = "🔻"
	}

	bigFloatValue := big.NewFloat(value)
	cv.Value = ac.FormatMoneyBigFloat(bigFloatValue)

	return cv
}
