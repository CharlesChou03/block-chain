package db

type Block struct {
	ID         int64   `gorm:"type:int(10) NOT NULL auto_increment;primary_key;" json:"id"`
	BlockNum   int64   `gorm:"uniqueIndex:block_num,type:int(10) NOT NULL;" json:"block_num"`
	BlockHash  string  `gorm:"type:varchar(64) NOT NULL;" json:"block_hash"`
	BlockTime  int64   `gorm:"type:int(13) NOT NULL;" json:"block_time"`
	ParentHash string  `gorm:"type:varchar(64) NOT NULL;" json:"parent_hash"`
	CreateTime []uint8 `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;" json:"create_time"`
	UpdateTime []uint8 `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;" json:"update_time"`
}

type Transaction struct {
	ID         int64   `gorm:"type:int(10) NOT NULL auto_increment;primary_key;" json:"id"`
	BlockId    int64   `gorm:"type:int(10) NOT NULL" json:"block_id"`
	TxHash     string  `gorm:"type:varchar(64) NOT NULL" json:"tx_hash"`
	From       string  `gorm:"type:varchar(64) DEFAULT NULL" json:"from"`
	To         string  `gorm:"type:varchar(64) DEFAULT NULL" json:"to"`
	Nonce      int64   `gorm:"type:int(10) DEFAULT NULL" json:"nonce"`
	Data       string  `gorm:"type:varchar(64) DEFAULT NULL" json:"data"`
	Value      string  `gorm:"type:varchar(64) DEFAULT NULL" json:"value"`
	CreateTime []uint8 `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;" json:"create_time"`
	UpdateTime []uint8 `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;" json:"update_time"`
	Block      Block   `gorm:"foreignKey:BlockId"`
}

type TransactionLog struct {
	ID          int64       `gorm:"type:int(10) NOT NULL auto_increment;primary_key;" json:"id"`
	TxId        int64       `gorm:"type:int(10) NOT NULL" json:"tx_id"`
	Index       int64       `gorm:"type:int(10) DEFAULT NULL" json:"index"`
	Data        string      `gorm:"type:varchar(64) DEFAULT NULL" json:"data"`
	CreateTime  []uint8     `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;" json:"create_time"`
	UpdateTime  []uint8     `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;" json:"update_time"`
	Transaction Transaction `gorm:"foreignKey:TxId"`
}
