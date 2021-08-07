package models

type GetTransactionByHashReq struct {
	TxHash string `uri:"txHash"`
}

func (r *GetTransactionByHashReq) Validate() bool {
	if r.TxHash == "" {
		return false
	}
	return true
}
