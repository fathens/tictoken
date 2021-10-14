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
	flag.Parse()
	args := flag.Args()
	hdpath := wallet.DefaultPath
	if len(args) < 1 {
		fmt.Println("No hdpath supplied. Use", hdpath)
	} else {
		hdpath = args[0]
	}

	file, err := ioutil.ReadFile("config.toml")
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

	account, err := seed.Derive(hdpath)
	if err != nil {
		panic(err)
	}
	address, err := account.Address()
	if err != nil {
		panic(err)
	}
	fmt.Println(address)
}
