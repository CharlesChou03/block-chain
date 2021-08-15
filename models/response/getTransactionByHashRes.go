package models

type GetTransactionByHashRes struct {
	TxHash string `json:"tx_hash"`
	From   string `json:"from"`
	To     string `json:"to"`
	Nonce  uint64 `json:"nonce"`
	Data   string `json:"data"`
	Value  string `json:"value"`
	Logs   []Log  `json:"logs"`
}

type Log struct {
	Index uint   `json:"index"`
	Data  string `json:"data"`
}
