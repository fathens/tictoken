package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fathens/tictoken/dapp"
	"github.com/fathens/tictoken/wallet"
	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	RpcServer string
}

func main() {
	configFile := flag.String("config", "config.toml", "Path of config")
	hdpath := flag.String("hdpath", wallet.DefaultPath, "HDPath")
	solc := flag.String("solc", "solc", "solc command")
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		panic("No filename supplied.")
	}
	fileName := args[0]
	cfg := readConfig(*configFile)
	fmt.Println("config =", cfg)

	mnemonic := os.Getenv("TICTOKEN_MNEMONIC")
	account := setupAccount(mnemonic, *hdpath)
	fmt.Println(account.Address())

	contractAddr, err := dapp.DeployFromSrc(account, *solc, fileName, "", "")
	if err != nil {
		panic(err)
	}
	fmt.Println(contractAddr)
}

func readConfig(path string) Config {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var cfg Config
	err = toml.Unmarshal(file, &cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

func setupAccount(mnemonic, hdpath string) wallet.Account {
	seed, err := wallet.InitByMnemonic(mnemonic)
	if err != nil {
		panic(err)
	}
	account, err := seed.Derive(hdpath)
	if err != nil {
		panic(err)
	}
	return account
}
