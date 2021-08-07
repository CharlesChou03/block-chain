package models

type GetBlockByNumReq struct {
	BlockNum int `uri:"id"`
}

func (r *GetBlockByNumReq) Validate() bool {
	if r.BlockNum < 0 {
		return false
	}
	return true
}
