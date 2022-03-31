package tests

import (
	mkputils "github.com/mkpproduction/mkp-sdk-go/mkp/utils"
	"log"
	"testing"
)

func TestCreateCredential(t *testing.T) {

	secret := "MKPMobile"
	value := "ultra voucher"

	result, err := mkputils.CreateCredential(secret, value)
	if err != nil {
		log.Println("TestCreateCredential Error:", err.Error())
		return
	}

	result2 := mkputils.Base64ToHex(result)

	log.Println("result:", result)
	log.Println("result2:", result2)
	log.Println("result2:", len(result2))

}
