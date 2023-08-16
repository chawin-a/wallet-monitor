package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/chawin-a/wallet-monitor/configs"
	"github.com/chawin-a/wallet-monitor/internal/explorer"
	"github.com/chawin-a/wallet-monitor/internal/postgres"
	"github.com/chawin-a/wallet-monitor/internal/worker"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	conf := configs.InitConfig()
	client, err := ethclient.Dial(conf.Node.Url)
	if err != nil {
		log.Panic(err)
	}
	defer client.Close()
	exp := explorer.NewExplorer(&conf.Explorer)
	pg := postgres.New(conf.Postgres)
	defer pg.Close()

	w := worker.NewWorker(exp, client, pg, conf.Worker)
	if err := w.Run(ctx); err != nil {
		log.Panic(err)
	}
	log.Println("gracefully shutdown")
}
