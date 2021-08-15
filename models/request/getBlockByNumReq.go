package models

type GetBlockByNumReq struct {
	BlockNum uint64 `uri:"id"`
}

func (r *GetBlockByNumReq) Validate() bool {
	if r.BlockNum < 0 {
		return false
	}
	return true
}
