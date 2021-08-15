package services

import (
	"time"

	"github.com/CharlesChou03/_git/block-chain.git/internal/db"
	"github.com/CharlesChou03/_git/block-chain.git/models"
	modelsReq "github.com/CharlesChou03/_git/block-chain.git/models/request"
	modelsRes "github.com/CharlesChou03/_git/block-chain.git/models/response"
)

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
	dbBlockData := db.Block{}
	if len(blockDataList) == 0 {
		return
	}
	date := time.Now().UTC()
	dbBlockData.CreateTime = date
	dbBlockData.UpdateTime = date

	blockTxMap := make(map[uint64][]string)
	for _, block := range blockDataList {
		dbBlockData.BlockNum = block.BlockNum
		dbBlockData.BlockHash = block.BlockHash
		dbBlockData.BlockTime = block.BlockTime
		dbBlockData.ParentHash = block.ParentHash
		dbBlockDataList = append(dbBlockDataList, dbBlockData)
		blockNumList = append(blockNumList, block.BlockNum)
		txList := block.Transactions
		for _, tx := range txList {
			blockTxMap[block.BlockNum] = append(blockTxMap[block.BlockNum], tx)
		}
	}

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

func getBlockFromEthereum(req *modelsReq.GetBlockByNumReq, res *modelsRes.GetBlockByNumRes) (int, models.BlockChainError) {
	statusCode, blockDataList, _ := GetBlockDataFromEthereum(0, req.BlockNum, false)
	if statusCode != 200 {
		return 424, models.FailedDependencyError
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

func getBlockFromDB(blockNum uint64, res *modelsRes.GetBlockByNumRes) (int, models.BlockChainError) {
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
		go saveTxDataListToDB(tx)
	}
	return 200, models.NoError
}

func saveTxDataListToDB(tx TransactionData) {
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
