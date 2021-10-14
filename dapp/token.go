package dapp

import "github.com/ethereum/go-ethereum/common/compiler"

func Compile(solcCmd, srcPath string) (map[string]*compiler.Contract, error) {
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
