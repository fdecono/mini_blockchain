package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

type Block struct {
	Index     int
	Timestamp string
	BPM       int
	Hash      string
	PrevHash  string
	Validator string
}

var Blockchain []Block
var bcServer chan []Block

var tempBlocks []Block
var candidateBlocks = make(chan Block)
var announcements = make(chan string)
var mutex = &sync.Mutex{}
var validators = make(map[string]int)

func calculateHash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	hashed := h.Sum(nil)

	return hex.EncodeToString(hashed)
}

func calculateBlockHash(block Block) string {
	record := string(block.Index) + block.Timestamp + string(block.BPM) + block.PrevHash
	return calculateHash(record)
}

func generateBlock(oldBlock Block, BPM int, address string) (Block, error) {
	var newBlock Block
	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Validator = address
	newBlock.Hash = calculateBlockHash(newBlock)

	return newBlock, nil
}

func isBlockValid(newBlock Block, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateBlockHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

// func run() error {
// 	mux := makeMuxRouter()
// 	httpAddr := os.Getenv("PORT")
// 	log.Println("Listening on ", httpAddr)
// 	s := &http.Server{
// 		Addr:           ":" + httpAddr,
// 		Handler:        mux,
// 		ReadTimeout:    10 * time.Second,
// 		WriteTimeout:   10 * time.Second,
// 		MaxHeaderBytes: 1 << 20,
// 	}

// 	if err := s.ListenAndServe(); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func makeMuxRouter() http.Handler {
// 	muxRouter := mux.NewRouter()
// 	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
// 	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
// 	return muxRouter
// }

// func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
// 	bytes, err := json.MarshalIndent(Blockchain, "", " ")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	io.WriteString(w, string(bytes))
// }

type Message struct {
	BPM int
}

// func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
// 	var m Message
// 	decoder := json.NewDecoder(r.Body)
// 	if err := decoder.Decode(&m); err != nil {
// 		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
// 		return
// 	}
// 	defer r.Body.Close()

// 	newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], m.BPM)
// 	if err != nil {
// 		respondWithJSON(w, r, http.StatusInternalServerError, err)
// 		return
// 	}
// 	if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
// 		newBlockchain := append(Blockchain, newBlock)
// 		replaceChain(newBlockchain)
// 		spew.Dump(Blockchain)
// 	}

// 	respondWithJSON(w, r, http.StatusCreated, newBlock)
// }

// func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload any) {
// 	response, err := json.MarshalIndent(payload, "", " ")
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte("HTTP 500: Internal Server Error"))
// 		return
// 	}
// 	w.WriteHeader(code)
// 	w.Write(response)
// }

func handleConn(conn net.Conn) {
	defer conn.Close()

	go func() {
		for {
			msg := <-announcements
			io.WriteString(conn, msg)
		}
	}()

	var address string

	io.WriteString(conn, "Enter token balance:")
	scanBalance := bufio.NewScanner(conn)
	for scanBalance.Scan() {
		balance, err := strconv.Atoi(scanBalance.Text())
		if err != nil {
			log.Printf("%v not a number: %v", scanBalance.Text(), err)
			return
		}
		t := time.Now()
		address = calculateHash(t.String())
		validators[address] = balance
		fmt.Println(validators)
		break
	}

	io.WriteString(conn, "\nEnter a new BPM:")
	scanBPM := bufio.NewScanner(conn)

	go func() {
		for {
			for scanBPM.Scan() {
				bpm, err := strconv.Atoi(scanBPM.Text())
				if err != nil {
					log.Printf("%v is not a number: %v", scanBPM.Text(), err)
					delete(validators, address)
					conn.Close()
				}
				mutex.Lock()
				oldLastIndex := Blockchain[len(Blockchain)-1]
				mutex.Unlock()
				newBlock, err := generateBlock(oldLastIndex, bpm, address)
				if err != nil {
					log.Println(err)
					continue
				}
				if isBlockValid(newBlock, oldLastIndex) {
					candidateBlocks <- newBlock
				}
				io.WriteString(conn, "\nEnter a new BPM:")
			}
		}
	}()

	for {
		time.Sleep(10 * time.Second)
		mutex.Lock()
		output, err := json.Marshal(Blockchain)
		mutex.Unlock()
		if err != nil {
			log.Fatal(err)
		}
		io.WriteString(conn, string(output)+"\n")
	}
}

func pickWinner() {
	time.Sleep(20 * time.Second)
	mutex.Lock()
	temp := tempBlocks
	mutex.Unlock()

	lotteryPool := []string{}
	if len(temp) > 0 {
	OUTER:
		for _, block := range temp {
			// if already in lotteryPool, skip:
			for _, node := range lotteryPool {
				if block.Validator == node {
					continue OUTER
				}
			}

			mutex.Lock()
			setValidators := validators
			mutex.Unlock()

			k, ok := setValidators[block.Validator]
			if ok {
				for i := 0; i < k; i++ {
					lotteryPool = append(lotteryPool, block.Validator)
				}
			}
		}

		s := rand.NewSource(time.Now().Unix())
		r := rand.New(s)
		lotteryWinner := lotteryPool[r.Intn(len(lotteryPool))]

		// add winner block to blockchain
		for _, block := range temp {
			if block.Validator == lotteryWinner {
				mutex.Lock()
				Blockchain = append(Blockchain, block)
				mutex.Unlock()
				for _ = range validators {
					announcements <- "\nWinning validator: " + lotteryWinner + "\n"
				}
				break
			}
		}
	}
	mutex.Lock()
	tempBlocks = []Block{}
	mutex.Unlock()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	t := time.Now()
	genesisBlock := Block{}
	genesisBlock = Block{0, t.String(), 0, calculateBlockHash(genesisBlock), "", ""}
	spew.Dump(genesisBlock)
	Blockchain = append(Blockchain, genesisBlock)

	// start TCP
	server, err := net.Listen("tcp", ":"+os.Getenv("TCP_PORT"))
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	go func() {
		for candidate := range candidateBlocks {
			mutex.Lock()
			tempBlocks = append(tempBlocks, candidate)
			mutex.Unlock()
			announcements <- "\nNew Block For Review\n"
		}
	}()

	go func() {
		for {
			pickWinner()
		}
	}()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
}
