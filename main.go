package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fathens/tictoken/wallet"
	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	RpcServer string
}

func main() {
	configFile := flag.String("config", "config.toml", "HDPath")
	hdpath := flag.String("hdpath", wallet.DefaultPath, "HDPath")
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		panic("No filename supplied.")
	}

	file, err := ioutil.ReadFile(*configFile)
	if err != nil {
		panic(err)
	}

	var cfg Config
	err = toml.Unmarshal(file, &cfg)
	if err != nil {
		panic(err)
	}

	mnemonic := os.Getenv("TICTOKEN_MNEMONIC")
	seed, err := wallet.InitByMnemonic(mnemonic)
	if err != nil {
		panic(err)
	}

	account, err := seed.Derive(*hdpath)
	if err != nil {
		panic(err)
	}
	address, err := account.Address()
	if err != nil {
		panic(err)
	}
	fmt.Println(address)
}
