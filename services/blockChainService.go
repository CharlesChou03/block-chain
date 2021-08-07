package services

import (
	"github.com/CharlesChou03/_git/block-chain.git/internal/db"
	"github.com/CharlesChou03/_git/block-chain.git/models"
	modelsReq "github.com/CharlesChou03/_git/block-chain.git/models/request"
	modelsRes "github.com/CharlesChou03/_git/block-chain.git/models/response"
)

func GetBlocks(req *modelsReq.GetLatestNBlockReq, res *modelsRes.GetLatestNBlockRes) (int, *modelsRes.GetLatestNBlockRes, models.BlockChainError) {
	req.Init()

	blocksFromDB := []db.Block{}
	db.MySQLDB.QueryBlocks(req.Limit, &blocksFromDB)
	if len(blocksFromDB) == 0 {
		return 204, res, models.NotFoundError
	}
	for _, block := range blocksFromDB {
		blockInfo := modelsRes.BlockInfo{}
		blockInfo.BlockNum = block.BlockNum
		blockInfo.BlockHash = block.BlockHash
		blockInfo.BlockTime = block.BlockTime
		blockInfo.ParentHash = block.ParentHash
		res.Blocks = append(res.Blocks, blockInfo)
	}
	return 200, res, models.NoError
}

func GetBlock(req *modelsReq.GetBlockByNumReq, res *modelsRes.GetBlockByNumRes) (int, *modelsRes.GetBlockByNumRes, models.BlockChainError) {
	if !req.Validate() {
		return 400, res, models.BadRequestError
	}
	blocksFromDB := []db.Block{}
	db.MySQLDB.QueryBlock(req.BlockNum, &blocksFromDB)
	if len(blocksFromDB) == 0 {
		return 204, res, models.NotFoundError
	}
	block := blocksFromDB[0]
	blockId := block.ID
	res.BlockNum = block.BlockNum
	res.BlockHash = block.BlockHash
	res.BlockTime = block.BlockTime
	res.ParentHash = block.ParentHash

	transactionsFromDB := []db.Transaction{}
	db.MySQLDB.QueryTransactionsByBlockId([]string{"tx_hash"}, blockId, &transactionsFromDB)

	if len(transactionsFromDB) == 0 {
		res.Transactions = []string{}
	} else {
		for _, transaction := range transactionsFromDB {
			res.Transactions = append(res.Transactions, transaction.TxHash)
		}
	}

	return 200, res, models.NoError
}

func GetTransaction(req *modelsReq.GetTransactionByHashReq, res *modelsRes.GetTransactionByHashRes) (int, *modelsRes.GetTransactionByHashRes, models.BlockChainError) {
	if !req.Validate() {
		return 400, res, models.BadRequestError
	}
	transactionsFromDB := []db.Transaction{}
	db.MySQLDB.QueryTransactionByTxHash([]string{"*"}, req.TxHash, &transactionsFromDB)
	if len(transactionsFromDB) == 0 {
		return 204, res, models.NotFoundError
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

	return 200, res, models.NoError
}
