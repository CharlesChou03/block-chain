package db

import "time"

type Block struct {
	ID         uint64    `gorm:"type:int(10) NOT NULL auto_increment;primary_key;" json:"id"`
	BlockNum   uint64    `gorm:"uniqueIndex:block_num,type:int(10) NOT NULL;" json:"block_num"`
	BlockHash  string    `gorm:"type:varchar(66) NOT NULL;" json:"block_hash"`
	BlockTime  uint64    `gorm:"type:int(13) NOT NULL;" json:"block_time"`
	ParentHash string    `gorm:"type:varchar(66) NOT NULL;" json:"parent_hash"`
	CreateTime time.Time `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;" json:"create_time"`
	UpdateTime time.Time `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;" json:"update_time"`
}

type BlockTransaction struct {
	ID      uint64 `gorm:"type:int(10) NOT NULL auto_increment;primary_key;" json:"id"`
	BlockId uint64 `gorm:"type:int(10) NOT NULL" json:"block_id"`
	TxHash  string `gorm:"type:varchar(66) NOT NULL" json:"tx_hash"`
	Block   Block  `gorm:"foreignKey:BlockId"`
}

type Transaction struct {
	ID         uint64    `gorm:"type:int(10) NOT NULL auto_increment;primary_key;" json:"id"`
	TxHash     string    `gorm:"type:varchar(66) NOT NULL" json:"tx_hash"`
	From       string    `gorm:"type:varchar(42) DEFAULT NULL" json:"from"`
	To         string    `gorm:"type:varchar(42) DEFAULT NULL" json:"to"`
	Nonce      uint64    `gorm:"type:int(10) DEFAULT NULL" json:"nonce"`
	Data       string    `gorm:"type:text DEFAULT NULL" json:"data"`
	Value      string    `gorm:"type:text DEFAULT NULL" json:"value"`
	CreateTime time.Time `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;" json:"create_time"`
	UpdateTime time.Time `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;" json:"update_time"`
}

type TransactionLog struct {
	ID          uint64      `gorm:"type:int(10) NOT NULL auto_increment;primary_key;" json:"id"`
	TxId        uint64      `gorm:"type:int(10) NOT NULL" json:"tx_id"`
	Index       uint        `gorm:"type:int(10) DEFAULT NULL" json:"index"`
	Data        string      `gorm:"type:text DEFAULT NULL" json:"data"`
	Transaction Transaction `gorm:"foreignKey:TxId"`
}
