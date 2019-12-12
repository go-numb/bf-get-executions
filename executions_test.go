package eoy

import (
	"os"
	"testing"

	"github.com/labstack/gommon/log"
)

func TestExecutions(t *testing.T) {
	key := os.Getenv("BFKEY")
	if key == "" {
		log.Fatal("undefine BFKEY")
	}
	secret := os.Getenv("BFSECRET")
	if key == "" {
		log.Fatal("undefine BFSECRET")
	}

	productCode := "FX_BTC_JPY"
	bf := New(key, secret, productCode)
	bf.Executions()
}
