package worker

import (
	"context"
	"log"
	"time"

	"github.com/chawin-a/wallet-monitor/internal/datagateway"
	"github.com/chawin-a/wallet-monitor/internal/entity"
	"github.com/chawin-a/wallet-monitor/internal/explorer"
	"github.com/chawin-a/wallet-monitor/internal/postgres"
	"github.com/chawin-a/wallet-monitor/internal/utils"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/sync/errgroup"
)

type Config struct {
	Id            int    `mapstructure:"id"`
	Interval      string `mapstructure:"interval"`
	MaxConcurrent int    `mapstructure:"max_concurrent"`
}

type Worker struct {
	id            int
	interval      time.Duration
	maxConcurrent int
	explorer      *explorer.Explorer
	ethcliet      *ethclient.Client
	walletTxDg    datagateway.WalletTransaction
	walletDg      datagateway.Wallet
}

func NewWorker(
	exp *explorer.Explorer,
	client *ethclient.Client,
	db *postgres.Postgres,
	conf Config,
) *Worker {
	return &Worker{
		id:            conf.Id,
		interval:      utils.Must(time.ParseDuration(conf.Interval)),
		maxConcurrent: conf.MaxConcurrent,
		explorer:      exp,
		ethcliet:      client,
		walletTxDg:    db,
		walletDg:      db,
	}
}

func (w *Worker) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.interval)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := w.Main(ctx); err != nil {
				return err
			}
		}
	}
}

func (w *Worker) Main(ctx context.Context) error {
	currentBlock, err := w.ethcliet.BlockNumber(ctx)
	if err != nil {
		return err
	}
	endBlock := currentBlock - 20
	wallet, err := w.walletDg.GetOldestLastBlockWalletByWorkerId(ctx, w.id)
	if err != nil {
		return err
	}
	if wallet == nil {
		return nil
	}
	response, err := w.explorer.TxList(ctx, wallet.Wallet, wallet.LastBlock, endBlock) // delay 20 blocks
	if err != nil {
		log.Println("something went wrong: " + err.Error()) // do not stop service
		return nil
	}
	if len(response.Result) == 0 {
		if _, err := w.walletDg.UpdateWalletLastBlock(ctx, wallet.Wallet, endBlock); err != nil {
			return err
		}
		return nil
	}
	txChan := make(chan *entity.Transaction, len(response.Result))
	errGroup, eCtx := errgroup.WithContext(ctx)
	for i := 0; i < w.maxConcurrent; i++ {
		errGroup.Go(func() error {
			for t := range txChan {
				walletTx := &entity.WalletTransaction{
					Wallet:           wallet.Wallet,
					BlockNumber:      t.BlockNumber,
					BlockHash:        t.BlockHash,
					From:             t.From,
					To:               t.To,
					Gas:              t.Gas,
					GasPrice:         t.GasPrice,
					Nonce:            t.Nonce,
					Hash:             t.Hash,
					Input:            t.Input,
					TransactionIndex: t.TransactionIndex,
					Value:            t.Value,
					Status:           t.Status,
					Type: func() entity.WalletTransactionType {
						if wallet.Wallet == t.From {
							return entity.OUTGOING
						} else if wallet.Wallet == t.To {
							return entity.INCOMING
						} else {
							return entity.UNKNOWN
						}
					}(),
				}
				if _, err := w.walletTxDg.UpsertWalletTransaction(eCtx, walletTx); err != nil {
					return err
				}
			}
			return nil
		})
	}

	lastBlock := uint64(0)
	errGroup.Go(func() error {
		for _, tx := range response.Result {
			lastBlock = max(lastBlock, tx.BlockNumber)
			txChan <- tx
		}
		close(txChan)
		return nil
	})

	if err := errGroup.Wait(); err != nil {
		return err
	}

	if _, err := w.walletDg.UpdateWalletLastBlock(ctx, wallet.Wallet, lastBlock); err != nil {
		return err
	}

	log.Println("process wallet", wallet.Wallet, "until block", lastBlock)

	return nil
}
