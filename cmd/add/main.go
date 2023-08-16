package main

import (
	"context"
	"log"

	"github.com/chawin-a/wallet-monitor/configs"
	"github.com/chawin-a/wallet-monitor/internal/entity"
	"github.com/chawin-a/wallet-monitor/internal/postgres"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	ctx := context.Background()
	conf := configs.InitConfig()
	pg := postgres.New(conf.Postgres)
	for _, wallet := range conf.Wallets {
		if _, err := pg.CreateWallet(ctx, &entity.Wallet{
			Wallet:    common.HexToAddress(wallet),
			WorkerId:  conf.Worker.Id,
			LastBlock: 0,
		}); err != nil {
			log.Panic(err)
		}
	}
}
