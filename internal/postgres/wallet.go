package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/chawin-a/wallet-monitor/internal/datagateway"
	"github.com/chawin-a/wallet-monitor/internal/entity"
	"github.com/ethereum/go-ethereum/common"
)

var _ datagateway.Wallet = (*Postgres)(nil)

type Wallet struct {
	Wallet    string
	WorkerId  int
	LastBlock uint64
}

func (w *Wallet) ToEntity() *entity.Wallet {
	return &entity.Wallet{
		Wallet:    common.HexToAddress(w.Wallet),
		WorkerId:  w.WorkerId,
		LastBlock: w.LastBlock,
	}
}

func fromWalletEntity(wallet *entity.Wallet) *Wallet {
	return &Wallet{
		Wallet:    strings.ToLower(wallet.Wallet.String()),
		WorkerId:  wallet.WorkerId,
		LastBlock: wallet.LastBlock,
	}
}

func (p *Postgres) GetOldestLastBlockWalletByWorkerId(ctx context.Context, workerId int) (*entity.Wallet, error) {
	wallet := &Wallet{}
	if err := p.DB.NewSelect().
		Model(wallet).
		Where("worker_id = ?", workerId).
		Order("last_block").
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return wallet.ToEntity(), nil
}

func (p *Postgres) UpdateWalletLastBlock(ctx context.Context, walletAddress common.Address, lastBlock uint64) (*entity.Wallet, error) {
	wallet := &Wallet{}
	if _, err := p.DB.NewUpdate().
		Model(wallet).
		Where("wallet = ?", strings.ToLower(walletAddress.String())).
		SetColumn("last_block", "?", lastBlock).
		Returning("*").
		Exec(ctx); err != nil {
		return nil, err
	}
	return wallet.ToEntity(), nil
}

func (p *Postgres) CreateWallet(ctx context.Context, wallet *entity.Wallet) (*entity.Wallet, error) {
	w := fromWalletEntity(wallet)
	if _, err := p.DB.NewInsert().
		Model(w).
		Returning("*").
		Exec(ctx); err != nil {
		return nil, err
	}
	return w.ToEntity(), nil
}
