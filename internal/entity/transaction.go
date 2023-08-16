package entity

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Transaction struct {
	BlockNumber      uint64
	BlockHash        common.Hash
	From             common.Address
	To               common.Address
	Gas              *big.Int
	GasPrice         *big.Int
	Nonce            uint64
	Hash             common.Hash
	Input            []byte
	TransactionIndex uint64
	Value            *big.Int
	Status           uint64
}
