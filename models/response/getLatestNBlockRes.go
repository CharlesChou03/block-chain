package models

type GetLatestNBlockRes struct {
	Blocks []BlockInfo `json:"blocks"`
}

type BlockInfo struct {
	BlockNum   uint64 `json:"block_num"`
	BlockHash  string `json:"block_hash"`
	BlockTime  uint64 `json:"block_time"`
	ParentHash string `json:"parent_hash"`
}
