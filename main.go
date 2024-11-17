package main

import (
	"log"

	"github.com/yourusername/blockchain_collector"
)

func main() {
	const (
		rpcURL     = "https://mainnet.infura.io/v3/YOUR_INFURA_PROJECT_ID"
		dbURL      = "user:password@tcp(localhost:3306)/blockchain" // Replace with your MySQL connection string
		batchSize  = 10
		blockRange = 1000
		startBlock = 17000000
		endBlock   = 17001000
	)

	collector, err := blockchain_collector.NewCollector(rpcURL, dbURL, batchSize, blockRange)
	if err != nil {
		log.Fatalf("Failed to initialize collector: %v\n", err)
	}
	defer collector.Close()

	if err := collector.CollectData(startBlock, endBlock); err != nil {
		log.Fatalf("Data collection failed: %v\n", err)
	}
}
