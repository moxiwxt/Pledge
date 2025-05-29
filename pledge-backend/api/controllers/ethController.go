package controllers

import (
	"github.com/gin-gonic/gin"
	"pledge-backend/api/common/statecode"
	"pledge-backend/api/models"
	"pledge-backend/api/models/request"
	"pledge-backend/api/models/response"
	"pledge-backend/api/services"
	"pledge-backend/api/validate"
)

type EthController struct {
}

func (c *EthController) GetBlockData(ctx *gin.Context) {

	res := response.Gin{Res: ctx}
	req := request.BlockData{}
	var result []models.BlockDataRes

	errCode := validate.NewBlockData().BlockData(ctx, &req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	errCode = services.NewGetBlockData().BlockData(req.BlockNumber, req.Full, &result)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	res.Response(ctx, statecode.CommonSuccess, result)
	return

}

func (c *EthController) SyncTransactionData(ctx *gin.Context) {

	res := response.Gin{Res: ctx}
	req := request.TxData{}
	var result []models.Transaction

	errCode := validate.NewTxData().TxData(ctx, &req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	errCode = services.NewTxData().TxData(req.TxHash, &result)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	res.Response(ctx, statecode.CommonSuccess, result)
	return

}
