package tests

import (
	"fmt"
	mkputils "github.com/mkpproduction/mkp-sdk-go/mkp/utils"
	"testing"
)

func TestCalcTaxAmount(t *testing.T) {

	amount := mkputils.CalcTaxAmount(mkputils.CalcAmount{
		Amount:    2869930.00,
		FlgIncTax: "Y",
		TaxPct:    11,
		Round:     2,
		RoundMode: "",
	})

	fmt.Println("Amount:", amount)
}
