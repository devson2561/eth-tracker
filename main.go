package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"sync"

	"github.com/devson2561/eth-tracker/db"
	"github.com/devson2561/eth-tracker/models"
	"github.com/devson2561/eth-tracker/repository"
	"github.com/devson2561/eth-tracker/utils"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

const (
	rpcUrl  = "wss://eth-sepolia.g.alchemy.com/v2/WEMf4aoFyj_A63aEKOzuXQyr0O9Uvdsa"
	dbFile  = "./ethereum-watcher.db"
	address = "0x28c6c06298d514db089934071355e5743bf21d60"
)

type App struct {
	repo            repository.Repository
	addressRepo     repository.AddressRepository
	transactionRepo repository.TransactionRepository
	blockRepo       repository.BlockRepository
	client          *ethclient.Client
	validAddresses  []string
}

func main() {
	app, err := NewApp(dbFile)
	if err != nil {
		logrus.Fatalf("Failed to initialize the app: %v", err)
	}
	app.run(os.Args)
}

func NewApp(dbFile string) (*App, error) {
	db, err := db.InitializeDatabase(dbFile)
	if err != nil {
		logrus.Fatalf("Failed to initialize the database: %v", err)
	}

	repo := repository.NewRepository(db)
	addressRepo := repository.NewAddressRepository(repo)
	transactionRepo := repository.NewTransactionRepository(repo)
	blockRepo := repository.NewBlockRepository(repo)

	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		logrus.Fatal(err)
	}

	app := &App{
		repo:            repo,
		addressRepo:     addressRepo,
		transactionRepo: transactionRepo,
		blockRepo:       blockRepo,
		client:          client,
	}

	return app, nil
}

func (app *App) run(args []string) {
	if len(args) > 2 && args[1] == "add" {
		addAddress(app, address)
	} else if args[1] == "start" {
		start(app)
	}
}

func start(app *App) {

	addresses, err := app.addressRepo.FindAll()
	if err != nil {
		logrus.Fatalf("Failed to query addresses: %v", err)
	}

	var addressStrings []string
	for _, address := range addresses {
		addressStrings = append(addressStrings, address.Address)
	}

	app.validAddresses = addressStrings

	syncPreviousTransactions(app)

	logrus.Info("========= START TX WATCHER ===============")

	headers := make(chan *types.Header)
	sub, err := app.client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		logrus.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			logrus.Fatal(err)
		case header := <-headers:
			block, err := app.client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				logrus.Fatal(err)
			}

			logrus.Info("New block: ", block.Number().Uint64())

			for _, tx := range block.Transactions() {
				processTransaction(app, tx, addressStrings)
			}

			insertBlock(app, block)
		}
	}
}

func addAddress(app *App, address string) {
	addressInDB, err := app.addressRepo.Find(address)
	if err != nil {
		logrus.Fatalf("Error: %v", err)
	}
	if addressInDB != nil {
		logrus.Fatalf("Failed to add the address: %s already exists", address)
	}

	err = app.addressRepo.Create(&models.Address{Address: address})
	if err != nil {
		logrus.Fatalf("Failed to add the address: %v", err)
	}
	logrus.Infof("Address %s added successfully", address)
}

func processTransaction(app *App, tx *types.Transaction, addresses []string) {
	var from string
	var to string

	if tx.To() != nil {
		to = tx.To().Hex()
	}

	from, err := utils.GetTxSender(tx)
	if err != nil {
		logrus.Warn("Failed to get sender address of tx:", tx.Hash().Hex())
		logrus.Error(err)
	}

	if utils.Contains(addresses, from) || utils.Contains(addresses, to) {
		txRaw, err := json.Marshal(tx)
		if err != nil {
			logrus.Errorf("Failed to read tx data: %v", err)
		}

		err = app.transactionRepo.Create(&models.Transaction{
			Hash:  tx.Hash().Hex(),
			Value: tx.Value().String(),
			TxRaw: string(txRaw),
			From:  from,
			To:    to,
		})
		if err != nil {
			logrus.Errorf("Failed to insert tx into the DB: %v", err)
		} else {
			logrus.Infof("Transaction %s inserted into the DB successfully.", tx.Hash().Hex())
		}

	}
}

func syncPreviousTransactions(app *App) {

	latestBlock, err := app.client.BlockByNumber(context.Background(), nil)
	if err != nil {
		logrus.Fatal(err)
	}

	lastBlockInDB, err := app.blockRepo.FindLatest()
	if err != nil || lastBlockInDB == nil {
		return
	}

	latestBlockNumber := latestBlock.Number().Uint64()

	if lastBlockInDB.BlockNumber < latestBlockNumber {
		logrus.Warnf("The DB is not up-to-date. Latest block in DB: %d, latest block in Ethereum: %d (%d blocks)",
			lastBlockInDB.BlockNumber, latestBlockNumber, latestBlockNumber-lastBlockInDB.BlockNumber)

		fmt.Println("Would you like to sync previous transactions? (yes/no)")
		var input string
		fmt.Scan(&input)

		input = strings.ToLower(input)

		if input == "yes" || input == "y" {
			processTransactionsByBlock(app, lastBlockInDB.BlockNumber+1, latestBlock.Number().Uint64())
		}

	} else {
		logrus.Info("The DB is up-to-date.")
	}
}

func processTransactionsByBlock(app *App, start uint64, end uint64) {
	var wg sync.WaitGroup
	for i := start; i <= end; i++ {
		wg.Add(1)
		go func(blockNumber uint64) {
			defer wg.Done()

			block, err := app.client.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
			if err != nil {
				log.Fatal(err)
			}

			for _, tx := range block.Transactions() {
				processTransaction(app, tx, app.validAddresses)
			}

			insertBlock(app, block)
			logrus.Infof("Sync block %d success.", block.Number().Uint64())

		}(i)
	}
	wg.Wait()
}

func insertBlock(app *App, block *types.Block) {
	err := app.blockRepo.Create(&models.Block{
		BlockNumber:      block.Number().Uint64(),
		BlockHash:        block.Hash().Hex(),
		TransactionCount: len(block.Transactions()),
	})

	if err != nil {
		logrus.Fatalf("Failed to insert block into the DB: %v", err)
	}
}
