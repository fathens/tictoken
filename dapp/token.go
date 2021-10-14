package dapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/compiler"
	"github.com/fathens/tictoken/wallet"
)

func DeployFromSrc(
	account wallet.Account,
	solcCmd, srcPath string,
	args []string,
) (common.Address, error) {
	empty := common.Address{}
	contracts, err := compileFile(solcCmd, srcPath)
	if err != nil {
		return empty, err
	}
	params := make([]interface{}, len(args))
	for i, s := range args {
		params[i] = s
	}
	for key, contract := range contracts {
		if strings.HasPrefix(key, srcPath) {
			return deployContract(account, contract, params...)
		}
	}
	return empty, nil
}

func compileFile(solcCmd, srcPath string) (map[string]*compiler.Contract, error) {
	solidity, err := compiler.SolidityVersion(solcCmd)
	if err != nil {
		return nil, err
	}
	contracts, err := compiler.CompileSolidity(solidity.Path, srcPath)
	if err != nil {
		return nil, err
	}
	return contracts, nil
}

func deployContract(
	account wallet.Account,
	contract *compiler.Contract,
	params ...interface{},
) (common.Address, error) {
	empty := common.Address{}
	rawAbi, err := json.Marshal(contract.Info.AbiDefinition)
	if err != nil {
		return empty, err
	}
	abi, err := abi.JSON(bytes.NewReader(rawAbi))
	if err != nil {
		return empty, err
	}

	code := common.FromHex(contract.Code)
	if err != nil {
		return empty, err
	}

	opts := bind.NewKeyedTransactor(&account.PrivateKey)

	contractAddr, tx, bound, err := bind.DeployContract(
		opts,
		abi,
		code,
		nil,
		params...,
	)
	if err != nil {
		return empty, err
	}
	fmt.Printf("Tx %v, Contract %v", tx, bound)

	return contractAddr, nil
}
