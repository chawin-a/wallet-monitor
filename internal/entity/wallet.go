package entity

import "github.com/ethereum/go-ethereum/common"

type Wallet struct {
	Wallet    common.Address
	WorkerId  int
	LastBlock uint64
}
