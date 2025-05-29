package request

type TxData struct {
	TxHash string `uri:"tx_hash" binding:"required"`
}
