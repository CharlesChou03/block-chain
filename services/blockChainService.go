package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/CharlesChou03/_git/block-chain.git/internal/db"
	"github.com/CharlesChou03/_git/block-chain.git/models"
	modelsReq "github.com/CharlesChou03/_git/block-chain.git/models/request"
	modelsRes "github.com/CharlesChou03/_git/block-chain.git/models/response"
)

var blockCacheKeyPrefix = "blockNum:"
var txCacheKeyPrefix = "txHash:"
var withoutExpired = int64(0)

func GetBlocks(req *modelsReq.GetLatestNBlockReq, res *modelsRes.GetLatestNBlockRes) (int, models.BlockChainError) {
	if !req.Validate() {
		return 400, models.BadRequestError
	}
	req.Init()

	return getBlockListFromEthereum(req, res)
}

func getBlockListFromEthereum(req *modelsReq.GetLatestNBlockReq, res *modelsRes.GetLatestNBlockRes) (int, models.BlockChainError) {
	statusCode, blockDataList, _ := GetBlockDataFromEthereum(req.Limit, 0, true)
	if statusCode != 200 {
		return 424, models.FailedDependencyError
	}
	if len(blockDataList) == 0 {
		return 204, models.NotFoundError
	}

	for _, block := range blockDataList {
		blockInfo := modelsRes.BlockInfo{}
		blockInfo.BlockNum = block.BlockNum
		blockInfo.BlockHash = block.BlockHash
		blockInfo.BlockTime = block.BlockTime
		blockInfo.ParentHash = block.ParentHash
		res.Blocks = append(res.Blocks, blockInfo)
	}

	go saveBlockDataToDB(blockDataList)
	return 200, models.NoError
}

func saveBlockDataToDB(blockDataList []BlockData) {
	dbBlockDataList := []db.Block{}
	blockNumList := []uint64{}
	if len(blockDataList) == 0 {
		return
	}

	date := time.Now().UTC()

	cacheKeyValuePair := make(map[string][]byte)

	blockTxMap := make(map[uint64][]string)
	for _, block := range blockDataList {
		cacheKey := blockCacheKeyPrefix + fmt.Sprint(block.BlockNum)
		cacheBlockData := genBlockCacheData(block)
		cacheValue, err := json.Marshal(cacheBlockData)
		if err == nil {
			cacheKeyValuePair[cacheKey] = cacheValue
		}

		dbBlockData := genBlockDBData(block)
		dbBlockData.CreateTime = date
		dbBlockData.UpdateTime = date
		dbBlockDataList = append(dbBlockDataList, dbBlockData)
		blockNumList = append(blockNumList, block.BlockNum)
		txList := block.Transactions
		for _, tx := range txList {
			blockTxMap[block.BlockNum] = append(blockTxMap[block.BlockNum], tx)
		}
	}

	db.RedisDB.MSetValueToCache(cacheKeyValuePair, withoutExpired)
	db.MySQLDB.InsertBlockDataList(dbBlockDataList)
	blocksFromDB := []db.Block{}
	queryResult := db.MySQLDB.QueryBlocksByNumList([]string{"id", "block_num"}, blockNumList, &blocksFromDB)
	if queryResult == true && len(blocksFromDB) != 0 {
		blockTxList := []db.BlockTransaction{}
		blockTx := db.BlockTransaction{}
		for _, blockFromDB := range blocksFromDB {
			id := blockFromDB.ID
			num := blockFromDB.BlockNum
			txs := blockTxMap[num]
			for _, tx := range txs {
				blockTx.BlockId = id
				blockTx.TxHash = tx
				blockTxList = append(blockTxList, blockTx)
			}
			if len(blockTxList) != 0 {
				db.MySQLDB.InsertBlockTxList(blockTxList)
			}
		}
	}
}

func genBlockCacheData(block BlockData) db.RedisBlockData {
	cacheBlockData := db.RedisBlockData{}
	cacheBlockData.BlockNum = block.BlockNum
	cacheBlockData.BlockHash = block.BlockHash
	cacheBlockData.BlockTime = block.BlockTime
	cacheBlockData.ParentHash = block.ParentHash
	cacheBlockData.Transactions = block.Transactions
	return cacheBlockData
}

func genBlockDBData(block BlockData) db.Block {
	dbBlockData := db.Block{}
	dbBlockData.BlockNum = block.BlockNum
	dbBlockData.BlockHash = block.BlockHash
	dbBlockData.BlockTime = block.BlockTime
	dbBlockData.ParentHash = block.ParentHash
	return dbBlockData
}

func GetBlock(req *modelsReq.GetBlockByNumReq, res *modelsRes.GetBlockByNumRes) (int, models.BlockChainError) {
	if !req.Validate() {
		return 400, models.BadRequestError
	}

	statusCode, _ := getBlockFromDB(req.BlockNum, res)
	if statusCode == 204 {
		return getBlockFromEthereum(req, res)
	}
	return 200, models.NoError
}

func getBlockFromDB(blockNum uint64, res *modelsRes.GetBlockByNumRes) (int, models.BlockChainError) {
	statusCode, _ := getBlockFromCache(blockNum, res)
	if statusCode == 200 {
		return 200, models.NoError
	}

	blocksFromDB := []db.Block{}
	db.MySQLDB.QueryBlock(blockNum, &blocksFromDB)
	if len(blocksFromDB) == 0 {
		return 204, models.NotFoundError
	}
	block := blocksFromDB[0]
	blockId := block.ID
	res.BlockNum = block.BlockNum
	res.BlockHash = block.BlockHash
	res.BlockTime = block.BlockTime
	res.ParentHash = block.ParentHash

	blockTxsFromDB := []db.BlockTransaction{}
	db.MySQLDB.QueryBlockTxsByBlockId([]string{"tx_hash"}, blockId, &blockTxsFromDB)

	if len(blockTxsFromDB) == 0 {
		res.Transactions = []string{}
	} else {
		for _, transaction := range blockTxsFromDB {
			res.Transactions = append(res.Transactions, transaction.TxHash)
		}
	}
	return 200, models.NoError
}

func getBlockFromCache(blockNum uint64, res *modelsRes.GetBlockByNumRes) (int, models.BlockChainError) {
	blockFromCache := db.RedisBlockData{}
	cacheKey := blockCacheKeyPrefix + fmt.Sprint(blockNum)
	cacheData := db.RedisDB.GetValueFromCache(cacheKey)
	if cacheData != "" {
		if err := json.Unmarshal([]byte(cacheData), &blockFromCache); err != nil {
			fmt.Printf("[getBlockFromCache] parse error: %+v", err)
		} else {
			res.BlockNum = blockFromCache.BlockNum
			res.BlockHash = blockFromCache.BlockHash
			res.BlockTime = blockFromCache.BlockTime
			res.ParentHash = blockFromCache.ParentHash
			res.Transactions = blockFromCache.Transactions
			return 200, models.NoError
		}
	}
	return 204, models.NotFoundError
}

func getBlockFromEthereum(req *modelsReq.GetBlockByNumReq, res *modelsRes.GetBlockByNumRes) (int, models.BlockChainError) {
	statusCode, blockDataList, _ := GetBlockDataFromEthereum(0, req.BlockNum, false)
	if statusCode != 200 {
		return 204, models.NotFoundError
	}
	if len(blockDataList) == 0 {
		return 204, models.NotFoundError
	}
	blockData := blockDataList[0]
	res.BlockNum = blockData.BlockNum
	res.BlockHash = blockData.BlockHash
	res.BlockTime = blockData.BlockTime
	res.ParentHash = blockData.ParentHash

	blockTxHashList := blockData.Transactions

	res.Transactions = []string{}
	for _, txHash := range blockTxHashList {
		res.Transactions = append(res.Transactions, txHash)
	}

	go saveBlockDataToDB(blockDataList)
	return 200, models.NoError
}

func GetTransaction(req *modelsReq.GetTransactionByHashReq, res *modelsRes.GetTransactionByHashRes) (int, models.BlockChainError) {
	if !req.Validate() {
		return 400, models.BadRequestError
	}

	statusCode, _ := getTxFromDB(req.TxHash, res)
	if statusCode == 204 {
		return getTxFromEthereum(req.TxHash, res)
	}

	return 200, models.NoError
}

func getTxFromDB(txHash string, res *modelsRes.GetTransactionByHashRes) (int, models.BlockChainError) {
	statusCode, _ := getTxFromCache(txHash, res)
	if statusCode == 200 {
		return 200, models.NoError
	}

	transactionsFromDB := []db.Transaction{}
	db.MySQLDB.QueryTransactionByTxHash([]string{"*"}, txHash, &transactionsFromDB)
	if len(transactionsFromDB) == 0 {
		return 204, models.NotFoundError
	}
	transaction := transactionsFromDB[0]
	txId := transaction.ID
	res.TxHash = transaction.TxHash
	res.From = transaction.From
	res.To = transaction.To
	res.Nonce = transaction.Nonce
	res.Data = transaction.Data
	res.Value = transaction.Value

	transactionLogsFromDB := []db.TransactionLog{}
	db.MySQLDB.QueryTransactionLogById([]string{"*"}, txId, &transactionLogsFromDB)

	if len(transactionLogsFromDB) == 0 {
		res.Logs = []modelsRes.Log{}
	} else {
		for _, txLog := range transactionLogsFromDB {
			log := modelsRes.Log{}
			log.Index = txLog.Index
			log.Data = txLog.Data
			res.Logs = append(res.Logs, log)
		}
	}
	return 200, models.NoError
}

func getTxFromCache(txHash string, res *modelsRes.GetTransactionByHashRes) (int, models.BlockChainError) {
	txFromCache := db.RedisTxData{}
	cacheKey := txCacheKeyPrefix + txHash
	cacheData := db.RedisDB.GetValueFromCache(cacheKey)
	if cacheData != "" {
		if err := json.Unmarshal([]byte(cacheData), &txFromCache); err != nil {
			fmt.Printf("[getTxFromCache] parse error: %+v", err)
		} else {
			res.TxHash = txFromCache.TxHash
			res.From = txFromCache.From
			res.To = txFromCache.To
			res.Nonce = txFromCache.Nonce
			res.Data = txFromCache.Data
			res.Value = txFromCache.Value
			if len(txFromCache.Logs) == 0 {
				res.Logs = []modelsRes.Log{}
			} else {
				for _, l := range txFromCache.Logs {
					log := modelsRes.Log{}
					log.Index = l.Index
					log.Data = l.Data
					res.Logs = append(res.Logs, log)
				}
			}
			return 200, models.NoError
		}
	}
	return 204, models.NotFoundError
}

func getTxFromEthereum(txHash string, res *modelsRes.GetTransactionByHashRes) (int, models.BlockChainError) {
	statusCode, tx, _ := GetTransactionFromEthereum(txHash)
	if statusCode != 200 {
		return statusCode, models.FailedDependencyError
	}
	res.TxHash = tx.TxHash
	res.From = tx.From
	res.To = tx.To
	res.Nonce = tx.Nonce
	res.Data = tx.Data
	res.Value = tx.Value

	if len(tx.Logs) != 0 {
		for _, l := range tx.Logs {
			log := modelsRes.Log{}
			log.Index = l.Index
			log.Data = l.Data
			res.Logs = append(res.Logs, log)
		}
	}
	if tx.IsPending == false {
		go saveTxDataToDB(tx)
	}
	return 200, models.NoError
}

func saveTxDataToDB(tx TransactionData) {
	cacheKey := txCacheKeyPrefix + tx.TxHash
	cacheTxData := genTxCacheData(tx)
	cacheValue, err := json.Marshal(cacheTxData)
	if err == nil {
		db.RedisDB.SetValueToCache(cacheKey, cacheValue, withoutExpired)
	}

	txDBData := genTxDBData(tx)
	date := time.Now().UTC()
	txDBData.CreateTime = date
	txDBData.UpdateTime = date
	db.MySQLDB.InsertTx(txDBData)

	transactionsFromDB := []db.Transaction{}
	queryResult := db.MySQLDB.QueryTransactionByTxHash([]string{"id"}, tx.TxHash, &transactionsFromDB)
	if queryResult == true && len(transactionsFromDB) != 0 {
		txId := transactionsFromDB[0].ID
		dbTxLogs := []db.TransactionLog{}
		if len(tx.Logs) != 0 {
			for _, l := range tx.Logs {
				log := db.TransactionLog{}
				log.TxId = txId
				log.Index = l.Index
				log.Data = l.Data
				dbTxLogs = append(dbTxLogs, log)
			}
			db.MySQLDB.InsertTxLogList(dbTxLogs)
		}
	}
}

func genTxCacheData(tx TransactionData) db.RedisTxData {
	cacheTxData := db.RedisTxData{}
	cacheTxData.TxHash = tx.TxHash
	cacheTxData.From = tx.From
	cacheTxData.To = tx.To
	cacheTxData.Nonce = tx.Nonce
	cacheTxData.Data = tx.Data
	cacheTxData.Value = tx.Value
	if len(tx.Logs) == 0 {
		cacheTxData.Logs = []db.RedisTxLogData{}
	} else {
		for _, l := range tx.Logs {
			log := db.RedisTxLogData{}
			log.Index = l.Index
			log.Data = l.Data
			cacheTxData.Logs = append(cacheTxData.Logs, log)
		}
	}
	return cacheTxData
}

func genTxDBData(tx TransactionData) db.Transaction {
	dbTx := db.Transaction{}
	dbTx.TxHash = tx.TxHash
	dbTx.From = tx.From
	dbTx.To = tx.To
	dbTx.Nonce = tx.Nonce
	dbTx.Data = tx.Data
	dbTx.Value = tx.Value
	return dbTx
}
