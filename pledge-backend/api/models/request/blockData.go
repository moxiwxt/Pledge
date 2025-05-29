package request

type BlockData struct {
	BlockNumber string `uri:"block_num" binding:"required"`
	Full        string `form:"full" default:"false"`
}
