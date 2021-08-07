package models

type GetLatestNBlockReq struct {
	Limit int `form:"limit"`
}

func (r *GetLatestNBlockReq) Init() {
	if r.Limit == 0 {
		r.Limit = 20
	}
}
