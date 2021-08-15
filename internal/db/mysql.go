package db

import (
	"fmt"

	"github.com/CharlesChou03/_git/block-chain.git/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MySQLBlockChainDB struct {
	DB *gorm.DB
}

var (
	MySQLDB                *MySQLBlockChainDB
	TableTutors            string = "tutors"
	TableTutorLessonPrices string = "tutor_lesson_prices"
	TableTutorLanguages    string = "tutor_languages"
	TableLanguages         string = "languages"
)

func SetupMySQLDB() *MySQLBlockChainDB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("connection to mysql failed:", err)
	}

	return &MySQLBlockChainDB{DB: db}
}

func (db *MySQLBlockChainDB) Close() {
	sqlDB, _ := db.DB.DB()
	sqlDB.Close()
}

func (db *MySQLBlockChainDB) CreateBlockChainTable() {
	db.DB.Set("gorm:table_options", "COLLATE=utf8mb4_general_ci").AutoMigrate(&BlockTransaction{})
	db.DB.Set("gorm:table_options", "COLLATE=utf8mb4_general_ci").AutoMigrate(&TransactionLog{})
}

func (db *MySQLBlockChainDB) QueryBlocks(size int, blocks *[]Block) bool {
	if result := db.DB.Order("block_time DESC").Limit(size).Find(&blocks); result.Error != nil {
		fmt.Print("+v%", result.Error)
		return false
	}
	return true
}

func (db *MySQLBlockChainDB) QueryBlock(blockNum uint64, blocks *[]Block) bool {
	if result := db.DB.Where("block_num = ?", blockNum).Find(&blocks); result.Error != nil {
		fmt.Print("+v%", result.Error)
		return false
	}
	return true
}

func (db *MySQLBlockChainDB) QueryBlocksByNumList(fields []string, blockNums []uint64, blocks *[]Block) bool {
	if result := db.DB.Select(fields).Where("block_num IN (?)", blockNums).Find(&blocks); result.Error != nil {
		fmt.Print("+v%", result.Error)
		return false
	}
	return true
}

func (db *MySQLBlockChainDB) QueryBlockTxsByBlockId(fields []string, blockId uint64, blockTxs *[]BlockTransaction) bool {
	if result := db.DB.Select(fields).Where("block_id = ?", blockId).Find(&blockTxs); result.Error != nil {
		fmt.Print("+v%", result.Error)
		return false
	}
	return true
}

func (db *MySQLBlockChainDB) QueryTransactionByTxHash(fields []string, txHash string, transactions *[]Transaction) bool {
	if result := db.DB.Select(fields).Where("tx_hash = ?", txHash).Find(&transactions); result.Error != nil {
		fmt.Print("+v%", result.Error)
		return false
	}
	return true
}

func (db *MySQLBlockChainDB) QueryTransactionLogById(fields []string, id uint64, transactionLogs *[]TransactionLog) bool {
	if result := db.DB.Select(fields).Where("tx_id = ?", id).Find(&transactionLogs); result.Error != nil {
		fmt.Print("+v%", result.Error)
		return false
	}
	return true
}

func (db *MySQLBlockChainDB) InsertBlockDataList(blockDataList []Block) {
	db.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "block_num"}},
		DoUpdates: clause.AssignmentColumns([]string{"block_hash", "block_time"}),
	}).Create(&blockDataList)
}

func (db *MySQLBlockChainDB) InsertBlockTxList(blockTxList []BlockTransaction) {
	db.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "block_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"tx_hash"}),
	}).Create(&blockTxList)
}

func (db *MySQLBlockChainDB) DeleteTxByHash(txHash string) {
	db.DB.Delete(Transaction{}, "tx_hash = ?", txHash)
}

func (db *MySQLBlockChainDB) InsertTx(tx Transaction) {
	db.DB.Create(&tx)
}

func (db *MySQLBlockChainDB) DeleteTxLogList(txId uint64) {
	db.DB.Delete(TransactionLog{}, "tx_id = ?", txId)
}

func (db *MySQLBlockChainDB) InsertTxLogList(txLogList []TransactionLog) {
	db.DB.Create(&txLogList)
}
