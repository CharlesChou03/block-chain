package services

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/CharlesChou03/_git/block-chain.git/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var one = int64(1)
var bigOne = big.NewInt(1)
var bigZero = big.NewInt(0)
var hexPrefix = "0x"

type BlockData struct {
	BlockNum     uint64
	BlockHash    string
	BlockTime    uint64
	ParentHash   string
	Transactions []string
}

type TransactionData struct {
	TxHash    string
	From      string
	To        string
	Nonce     uint64
	Data      string
	Value     string
	Logs      []TransactionLog
	IsPending bool
}

type TransactionLog struct {
	Index uint
	Data  string
}

func GetBlockDataFromEthereum(limit int, blockNum uint64, isRange bool) (int, []BlockData, error) {
	if isRange == true {
		return getBlockList(limit)
	} else {
		return getSpecificBlock(blockNum)
	}
}

func getSpecificBlock(blockNum uint64) (int, []BlockData, error) {
	client, err := ethclient.Dial(config.RPCEndpoint)
	defer client.Close()
	blockDataList := []BlockData{}
	if err != nil {
		fmt.Printf("%+v\n", err)
		return 424, blockDataList, err
	}
	specificBlock, err := client.BlockByNumber(context.Background(), new(big.Int).SetUint64(blockNum))
	if err != nil {
		fmt.Printf("%+v\n", err)
		return 424, blockDataList, err
	}
	specificBlockData := genBlockData(specificBlock)
	blockDataList = append(blockDataList, specificBlockData)
	return 200, blockDataList, nil
}

func getBlockList(limit int) (int, []BlockData, error) {
	client, err := ethclient.Dial(config.RPCEndpoint)
	defer client.Close()
	blockDataList := []BlockData{}
	if err != nil {
		fmt.Printf("%+v\n", err)
		return 424, blockDataList, err
	}
	latestBlock, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return 424, blockDataList, err
	}
	latestBlockNum := latestBlock.Header().Number
	start := getStartNumber(limit, latestBlockNum)

	blockChannel := make(chan *types.Block)
	channelNum := 0
	for i := new(big.Int).Set(start); i.Cmp(latestBlockNum) < 0; i.Add(i, bigOne) {
		channelNum += 1
		newI := new(big.Int).Set(i)
		go func(v big.Int) {
			block, _ := client.BlockByNumber(context.Background(), &v)
			blockChannel <- block
		}(*newI)
	}
	blockDataList = genBlockDataList(blockChannel, channelNum, latestBlock)
	return 200, blockDataList, nil
}

func getStartNumber(limit int, latestNum *big.Int) *big.Int {
	num := big.NewInt(int64(limit) - one)
	start := new(big.Int).Set(latestNum)
	start.Sub(start, num)

	if start.Cmp(bigZero) < 0 {
		start = big.NewInt(0)
	}
	return start
}

func genBlockDataList(blockChannel chan *types.Block, channelNum int, latestBlock *types.Block) []BlockData {
	blockDataList := []BlockData{}

	latestBlockData := genBlockData(latestBlock)
	blockDataList = append(blockDataList, latestBlockData)
	for i := 0; i < channelNum; i++ {
		block := <-blockChannel
		blockData := genBlockData(block)
		blockDataList = append(blockDataList, blockData)
	}
	return blockDataList
}

func genBlockData(block *types.Block) BlockData {
	blockData := BlockData{}
	blockData.BlockNum = block.NumberU64()
	blockData.BlockHash = block.Hash().Hex()
	blockData.BlockTime = block.Time()
	blockData.ParentHash = block.ParentHash().Hex()

	txDataList := []string{}
	txs := block.Transactions()
	for j := 0; j < len(txs); j++ {
		tx := txs[j]
		txDataList = append(txDataList, tx.Hash().Hex())
	}
	blockData.Transactions = txDataList
	return blockData
}

func GetTransactionFromEthereum(txHashHex string) (int, TransactionData, error) {
	client, err := ethclient.Dial(config.RPCEndpoint)
	defer client.Close()
	transactionData := TransactionData{}
	if err != nil {
		fmt.Printf("ethclient error: %+v\n", err)
		return 424, transactionData, err
	}

	txHash := common.HexToHash(txHashHex)
	tx, isPending, err := client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		fmt.Printf("TransactionByHash error: %+v\n", err)
		return 204, transactionData, err
	}

	msg, err := tx.AsMessage(types.NewEIP155Signer(tx.ChainId()), nil)
	if err != nil {
		fmt.Printf("AsMessage error: %+v\n", err)
		return 424, transactionData, err
	}

	transactionData.TxHash = tx.Hash().Hex()
	transactionData.From = msg.From().Hex()
	transactionData.To = msg.To().Hex()
	transactionData.Nonce = tx.Nonce()
	transactionData.Data = hexPrefix + hex.EncodeToString(tx.Data())
	transactionData.Value = msg.Value().String()
	transactionData.IsPending = isPending
	txReceipt, _ := client.TransactionReceipt(context.Background(), txHash)

	txLogs := []TransactionLog{}
	for _, log := range txReceipt.Logs {
		txLog := TransactionLog{}
		txLog.Index = log.Index
		txLog.Data = hexPrefix + hex.EncodeToString(log.Data)
		txLogs = append(txLogs, txLog)
	}
	transactionData.Logs = txLogs

	return 200, transactionData, nil
}
