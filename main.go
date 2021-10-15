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
	PrivateKey string
}

func main() {
	configFile := flag.String("config", "config.toml", "Path of config")
	hdpath := flag.String("hdpath", wallet.DefaultPath, "HDPath")
	solc := flag.String("solc", "solc", "solc command")
	flag.Parse()
	fullArgs := flag.Args()
	if len(fullArgs) < 1 {
		panic("command name (deploy or invoke) must be supplied.")
	}
	cmd := fullArgs[0]
	args := fullArgs[1:]

	cfg := readConfig(*configFile)
	fmt.Println("config:", cfg)

	var account wallet.Account
	if len(cfg.PrivateKey) == 0 {
		mnemonic := os.Getenv("TICTOKEN_MNEMONIC")
		account = setupAccount(mnemonic, *hdpath)
	} else {
		a, err := wallet.ReadPrivateKey(cfg.PrivateKey)
		if err != nil {
			panic(err)
		}
		account = *a
	}
	fmt.Println("account:", account.Address())

	switch cmd {
	case "deploy": deploy(cfg, account, *solc, args)
	case "invoke": invoke(cfg, account, args)
	default: panic(fmt.Sprintf("Unsupported command: %v", cmd))
	}
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
	return *account
}

func deploy(config Config, account wallet.Account, solc string, args []string) {
	fmt.Println("Exec deploy command: ", args)
	if len(args) < 1 {
		panic("filename must be supplied.")
	}
	fileName := args[0]

	contractAddr, err := dapp.DeployFromSrc(config.RpcServer, account, solc, fileName, args[1:])
	if err != nil {
		panic(err)
	}
	fmt.Println(contractAddr)
}

func invoke(config Config, account wallet.Account, args []string) {
	fmt.Println("Exec invoke command: ", args)
}
