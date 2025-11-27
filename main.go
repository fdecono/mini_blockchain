package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Block struct {
	Index     int
	Timestamp string
	BPM       int
	Hash      string
	PrevHash  string
}

var Blockchain []Block
var BlockchainMutex sync.Mutex

func calculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + string(block.BPM) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)

	return hex.EncodeToString(hashed)
}

func generateBlock(oldBlock Block, BPM int) (Block, error) {
	var newBlock Block
	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock, nil
}

func isBlockValid(newBlock Block, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

func replaceChain(newBlocks []Block) {
	BlockchainMutex.Lock()
	defer BlockchainMutex.Unlock()
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

func run() error {
	mux := makeMuxRouter()
	httpAddr := os.Getenv("PORT")
	log.Println("Listening on ", httpAddr)
	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	muxRouter.HandleFunc("/stream", handleStreamBlockchain).Methods("GET")
	// Serve the HTML file for live viewing
	muxRouter.HandleFunc("/view", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	return muxRouter
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	BlockchainMutex.Lock()
	bytes, err := json.MarshalIndent(Blockchain, "", " ")
	BlockchainMutex.Unlock()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

type Message struct {
	BPM int
}

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	var m Message
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	BlockchainMutex.Lock()
	lastBlock := Blockchain[len(Blockchain)-1]
	BlockchainMutex.Unlock()

	newBlock, err := generateBlock(lastBlock, m.BPM)
	if err != nil {
		respondWithJSON(w, r, http.StatusInternalServerError, err)
		return
	}
	if isBlockValid(newBlock, lastBlock) {
		BlockchainMutex.Lock()
		newBlockchain := append(Blockchain, newBlock)
		replaceChain(newBlockchain)
		BlockchainMutex.Unlock()
		spew.Dump(Blockchain)
	}

	respondWithJSON(w, r, http.StatusCreated, newBlock)
}

// handleStreamBlockchain handles Server-Sent Events (SSE) for live blockchain updates
func handleStreamBlockchain(w http.ResponseWriter, r *http.Request) {
	log.Println("SSE client connected from:", r.RemoteAddr)
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Get flusher to ensure we can flush data
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Create a ticker to send updates every second
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// Send initial blockchain state
	log.Println("Sending initial blockchain state")
	sendBlockchainUpdate(w)
	flusher.Flush()

	// Keep connection alive and send updates
	for {
		select {
		case <-ticker.C:
			sendBlockchainUpdate(w)
			flusher.Flush()
		case <-r.Context().Done():
			// Client disconnected
			log.Println("SSE client disconnected")
			return
		}
	}
}

// sendBlockchainUpdate sends the current blockchain state as an SSE event
func sendBlockchainUpdate(w http.ResponseWriter) {
	BlockchainMutex.Lock()
	blockCount := len(Blockchain)
	// Use Marshal without indent to avoid newline issues in SSE
	bytes, err := json.Marshal(Blockchain)
	BlockchainMutex.Unlock()

	if err != nil {
		log.Printf("Error marshaling blockchain: %v", err)
		fmt.Fprintf(w, "event: error\ndata: %s\n\n", err.Error())
		return
	}

	// Send as SSE event - single line JSON works best with SSE
	fmt.Fprintf(w, "data: %s\n\n", string(bytes))
	log.Printf("Sent blockchain update (%d blocks) via SSE", blockCount)
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload any) {
	response, err := json.MarshalIndent(payload, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize genesis block synchronously
	t := time.Now()
	genesisBlock := Block{0, t.String(), 0, "", ""}
	// Calculate hash for genesis block
	genesisBlock.Hash = calculateHash(genesisBlock)
	spew.Dump(genesisBlock)
	BlockchainMutex.Lock()
	Blockchain = append(Blockchain, genesisBlock)
	BlockchainMutex.Unlock()

	// Start the block generator
	startBlockGenerator()

	log.Fatal(run())
}
