package datagateway

import (
	"context"

	"github.com/chawin-a/wallet-monitor/internal/entity"
)

type WalletTransaction interface {
	UpsertWalletTransaction(ctx context.Context, tx *entity.WalletTransaction) (*entity.WalletTransaction, error)
}
