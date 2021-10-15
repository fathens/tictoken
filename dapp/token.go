package dapp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/compiler"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fathens/tictoken/wallet"
)

func DeployFromSrc(
	rpcserver string,
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
			return deployContract(rpcserver, account, contract, params...)
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
	rpcserver string,
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

	client, err := ethclient.Dial(rpcserver)
	if err != nil {
		return empty, err
	}
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return empty, err
	}

	opts, err := bind.NewKeyedTransactorWithChainID(&account.PrivateKey, chainId)
	if err != nil {
		return empty, err
	}

	contractAddr, tx, bound, err := bind.DeployContract(
		opts,
		abi,
		code,
		client,
		params...,
	)
	if err != nil {
		return empty, err
	}
	fmt.Printf("Tx %v, Contract %v", tx, bound)

	return contractAddr, nil
}
