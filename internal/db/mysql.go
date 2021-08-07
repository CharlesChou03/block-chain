package db

import (
	"fmt"

	"github.com/CharlesChou03/_git/block-chain.git/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
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
	db.DB.Set("gorm:table_options", "COLLATE=utf8mb4_general_ci").AutoMigrate(&TransactionLog{})
}

func (db *MySQLBlockChainDB) QueryBlocks(size int, blocks *[]Block) bool {
	if result := db.DB.Order("block_time DESC").Limit(size).Find(&blocks); result.Error != nil {
		fmt.Print("+v%", result.Error)
		return false
	}
	return true
}

func (db *MySQLBlockChainDB) QueryBlock(blockNum int, blocks *[]Block) bool {
	if result := db.DB.Where("block_num = ?", blockNum).Find(&blocks); result.Error != nil {
		fmt.Print("+v%", result.Error)
		return false
	}
	return true
}

func (db *MySQLBlockChainDB) QueryTransactionsByBlockId(fields []string, blockId int64, transactions *[]Transaction) bool {
	if result := db.DB.Select(fields).Where("block_id = ?", blockId).Find(&transactions); result.Error != nil {
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

func (db *MySQLBlockChainDB) QueryTransactionLogById(fields []string, id int64, transactionLogs *[]TransactionLog) bool {
	if result := db.DB.Select(fields).Where("tx_id = ?", id).Find(&transactionLogs); result.Error != nil {
		fmt.Print("+v%", result.Error)
		return false
	}
	return true
}
