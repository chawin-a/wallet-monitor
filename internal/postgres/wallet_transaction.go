package postgres

import (
	"context"
	"math/big"
	"strings"

	"github.com/chawin-a/wallet-monitor/internal/datagateway"
	"github.com/chawin-a/wallet-monitor/internal/entity"
	"github.com/chawin-a/wallet-monitor/internal/utils"
	"github.com/ethereum/go-ethereum/common"
)

var _ datagateway.WalletTransaction = (*Postgres)(nil)

type WalletTransaction struct {
	Wallet           string `bun:"wallet,pk"`
	Hash             string `bun:"tx_hash,pk"`
	BlockNumber      uint64 `bun:"block_number"`
	BlockHash        string `bun:"block_hash"`
	From             string `bun:"from"`
	To               string `bun:"to"`
	Gas              string `bun:"gas"`
	GasPrice         string `bun:"gas_price"`
	Nonce            uint64 `bun:"nonce"`
	Input            []byte `bun:"input"`
	TransactionIndex uint64 `bun:"transaction_index"`
	Value            string `bun:"value"`
	Status           uint64 `bun:"status"`
	Type             string `bun:"type"`
}

func fromWalletTransactionEntity(w *entity.WalletTransaction) *WalletTransaction {
	return &WalletTransaction{
		Wallet:           strings.ToLower(w.Wallet.String()),
		Hash:             strings.ToLower(w.Hash.String()),
		BlockNumber:      w.BlockNumber,
		BlockHash:        strings.ToLower(w.BlockHash.String()),
		From:             strings.ToLower(w.From.String()),
		To:               strings.ToLower(w.To.String()),
		Gas:              w.Gas.String(),
		GasPrice:         w.GasPrice.String(),
		Nonce:            w.Nonce,
		Input:            w.Input,
		TransactionIndex: w.TransactionIndex,
		Value:            w.Value.String(),
		Status:           w.Status,
		Type:             string(w.Type),
	}
}

func (w *WalletTransaction) ToEntity() *entity.WalletTransaction {
	return &entity.WalletTransaction{
		Wallet:           common.HexToAddress(w.Wallet),
		BlockNumber:      w.BlockNumber,
		BlockHash:        common.HexToHash(w.BlockHash),
		From:             common.HexToAddress(w.From),
		To:               common.HexToAddress(w.To),
		Gas:              utils.MustOk(new(big.Int).SetString(w.Gas, 10)),
		GasPrice:         utils.MustOk(new(big.Int).SetString(w.GasPrice, 10)),
		Nonce:            w.Nonce,
		Hash:             common.HexToHash(w.Hash),
		Input:            w.Input,
		TransactionIndex: w.TransactionIndex,
		Value:            utils.MustOk(new(big.Int).SetString(w.Value, 10)),
		Status:           w.Status,
		Type:             entity.WalletTransactionType(w.Type),
	}
}

// Upsert implements datagateway.WalletTransaction.
func (p *Postgres) UpsertWalletTransaction(ctx context.Context, tx *entity.WalletTransaction) (*entity.WalletTransaction, error) {
	walletTx := fromWalletTransactionEntity(tx)
	if _, err := p.DB.NewInsert().
		Model(walletTx).
		Ignore().
		Exec(ctx); err != nil {
		return nil, err
	}
	return walletTx.ToEntity(), nil
}
