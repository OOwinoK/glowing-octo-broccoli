package blockchain_collector

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Collector represents the blockchain data collector
type Collector struct {
	Client     *ethclient.Client
	DBPool     *pgxpool.Pool
	BatchSize  int
	BlockRange int64
}

// NewCollector initializes a new Collector instance
func NewCollector(rpcURL, dbURL string, batchSize int, blockRange int64) (*Collector, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}

	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, err
	}

	return &Collector{
		Client:     client,
		DBPool:     dbPool,
		BatchSize:  batchSize,
		BlockRange: blockRange,
	}, nil
}

// CollectData starts collecting blockchain data
func (c *Collector) CollectData(startBlock, endBlock int64) error {
	var wg sync.WaitGroup
	blockCh := make(chan int64, c.BlockRange)

	// Worker goroutines to fetch and process blocks
	for i := 0; i < c.BatchSize; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for blockNum := range blockCh {
				if err := c.processBlock(blockNum); err != nil {
					log.Printf("Error processing block %d: %v\n", blockNum, err)
				}
			}
		}()
	}

	// Push block numbers to the channel
	go func() {
		for block := startBlock; block <= endBlock; block++ {
			blockCh <- block
		}
		close(blockCh)
	}()

	wg.Wait()
	return nil
}

// processBlock fetches and stores data for a single block
func (c *Collector) processBlock(blockNumber int64) error {
	block, err := c.Client.BlockByNumber(context.Background(), ethereum.BigInt(blockNumber))
	if err != nil {
		return err
	}

	// Example: store block hash and number
	query := `INSERT INTO blocks (number, hash, timestamp) VALUES ($1, $2, $3)`
	_, err = c.DBPool.Exec(context.Background(), query, block.Number().Int64(), block.Hash().Hex(), block.Time())
	if err != nil {
		return err
	}

	log.Printf("Block %d stored successfully\n", blockNumber)
	return nil
}

// Close releases resources used by the collector
func (c *Collector) Close() {
	c.DBPool.Close()
}
