package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/davecgh/go-spew/spew"
)

// startBlockGenerator starts a goroutine that generates new blocks every 2 seconds
// with random BPM values between 60 and 100
func startBlockGenerator() {
	log.Println("Block generator started - will create blocks every 2 seconds")
	ticker := time.NewTicker(2 * time.Second)
	go func() {
		for range ticker.C {
			// Generate random BPM between 60 and 100
			randomBPM := rand.Intn(41) + 60 // 60-100 range

			// Get the last block safely
			BlockchainMutex.Lock()
			if len(Blockchain) == 0 {
				BlockchainMutex.Unlock()
				continue
			}
			lastBlock := Blockchain[len(Blockchain)-1]
			BlockchainMutex.Unlock()

			// Generate new block
			newBlock, err := generateBlock(lastBlock, randomBPM)
			if err != nil {
				log.Printf("Error generating block: %v", err)
				continue
			}

			// Validate the block
			if isBlockValid(newBlock, lastBlock) {
				// Use mutex to safely append to blockchain
				BlockchainMutex.Lock()
				Blockchain = append(Blockchain, newBlock)
				BlockchainMutex.Unlock()

				spew.Dump(newBlock)
				log.Printf("New block generated: Index=%d, BPM=%d, Hash=%s", newBlock.Index, newBlock.BPM, newBlock.Hash)
			} else {
				log.Printf("Block validation failed for block Index=%d", newBlock.Index)
			}
		}
	}()
}
