package wallet

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

const DefaultPath = "m/44'/60'/0'/0/0"

type Seed struct {
	bytes []byte
}

type Account struct {
	PrivateKey ecdsa.PrivateKey
}

func InitByMnemonic(words string) (Seed, error) {
	seed, err := bip39.NewSeedWithErrorChecking(words, "")
	return Seed{seed}, err
}

func (seed Seed) Derive(path string) (*Account, error) {
	hd, err := accounts.ParseDerivationPath(path)
	if err != nil {
		return nil, err
	}
	key, err := bip32.NewMasterKey(seed.bytes)
	if err != nil {
		return nil, err
	}

	for _, i := range hd {
		key, err = key.NewChildKey(i)
		if err != nil {
			return nil, err
		}
	}
	prv, err := crypto.ToECDSA(key.Key)
	if err != nil {
		return nil, err
	}

	return &Account{*prv}, nil
}

func ReadPrivateKey(hex string) (*Account, error) {
	key, err := crypto.HexToECDSA(hex)
	if err != nil {
		return nil, err
	}
	return &Account{*key}, nil
}

func (a Account) PublicKey() ecdsa.PublicKey {
	return a.PrivateKey.PublicKey
}

func (a Account) Address() common.Address {
	return crypto.PubkeyToAddress(a.PublicKey())
}
