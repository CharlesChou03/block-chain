package models

type GetLatestNBlockReq struct {
	Limit int `form:"limit"`
}

func (r *GetLatestNBlockReq) Validate() bool {
	if r.Limit < 0 || r.Limit > 30 {
		return false
	}
	return true
}

func (r *GetLatestNBlockReq) Init() {
	if r.Limit == 0 {
		r.Limit = 20
	}
}
