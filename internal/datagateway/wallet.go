package datagateway

import (
	"context"

	"github.com/chawin-a/wallet-monitor/internal/entity"
	"github.com/ethereum/go-ethereum/common"
)

type Wallet interface {
	GetOldestLastBlockWalletByWorkerId(ctx context.Context, WorkerId int) (*entity.Wallet, error)
	UpdateWalletLastBlock(ctx context.Context, wallet common.Address, lastBlock uint64) (*entity.Wallet, error)
	CreateWallet(ctx context.Context, wallet *entity.Wallet) (*entity.Wallet, error)
}
