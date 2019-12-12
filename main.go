package eoy

import (
	"flag"
	"os"

	"github.com/labstack/gommon/log"
)

var code string

func init() {
	flag.StringVar(&code, "code", "BTC_JPY", "use -code set product_code, default code is BTC_JPY")
	flag.Parse()
}

func main() {
	key := os.Getenv("BFKEY")
	if key == "" {
		log.Fatal("undefine BFKEY")
	}
	secret := os.Getenv("BFSECRET")
	if key == "" {
		log.Fatal("undefine BFSECRET")
	}

	bf := New(key, secret, code)
	bf.Executions()
}
