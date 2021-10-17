package dapp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
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
) (*common.Address, error) {
	contracts, err := compileFile(solcCmd, srcPath)
	if err != nil {
		return nil, err
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
	panic("Compiled file must exist.")
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
) (*common.Address, error) {
	rawAbi, err := json.Marshal(contract.Info.AbiDefinition)
	if err != nil {
		return nil, err
	}
	abi, err := abi.JSON(bytes.NewReader(rawAbi))
	if err != nil {
		return nil, err
	}

	code := common.FromHex(contract.Code)
	if err != nil {
		return nil, err
	}

	client, err := ethclient.Dial(rpcserver)
	if err != nil {
		return nil, err
	}
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	opts, err := bind.NewKeyedTransactorWithChainID(&account.PrivateKey, chainId)
	if err != nil {
		return nil, err
	}

	contractAddr, tx, bound, err := bind.DeployContract(
		opts,
		abi,
		code,
		client,
		params...,
	)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Tx %v, Contract %v\n", tx, bound)

	return &contractAddr, nil
}

func Invoke(
	rpcserver string,
	account wallet.Account,
	contractAddress common.Address,
	method string,
	abiJsonFile string,
	args []string,
) (*[]interface{}, error) {
	reader, err := os.Open(abiJsonFile)
	if err != nil {
		return nil, err
	}
	abiInfo, err := abi.JSON(reader)
	if err != nil {
		return nil, err
	}

	m := abiInfo.Methods[method]

	params := make([]interface{}, len(args))
	for i, s := range args {
		switch m.Inputs[i].Type.T {
		case abi.AddressTy:
			params[i] = common.HexToAddress(s)
		case abi.StringTy:
			params[i] = s
		default:
			params[i] = s
		}
	}

	return invokeContract(rpcserver, account.Address(), contractAddress, method, abiInfo, params...)
}

func invokeContract(
	rpcserver string,
	account, address common.Address,
	method string,
	abiInfo abi.ABI,
	params ...interface{},
) (*[]interface{}, error) {
	client, err := ethclient.Dial(rpcserver)
	if err != nil {
		return nil, err
	}

	opts := bind.CallOpts{
		From:    account,
		Context: context.Background(),
	}
	fmt.Println("Client:", client, opts, params)

	c := bind.NewBoundContract(address, abiInfo, client, client, client)

	results := make([]interface{}, 0)
	err = c.Call(&opts, &results, method, params...)
	if err != nil {
		return nil, err
	}
	return &results, nil
}
