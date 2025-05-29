package services

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"pledge-backend/api/common/statecode"
	"pledge-backend/api/models"
	"pledge-backend/db"
	"pledge-backend/log"
)

type TxDataService struct{}

func NewTxData() *TxDataService { return &TxDataService{} }

func (s *TxDataService) TxData(txHash string, result *[]models.Transaction) int {

	err := models.NewTxData().TxData(txHash, result)
	if err != nil {
		log.Logger.Error(err.Error())
		return statecode.CommonErrServerErr
	}

	client, err := ethclient.Dial("https://mainnet.infura.io/v3/YOUR_API_KEY")
	if err != nil {
		fmt.Println("连接失败:", err)
		return statecode.CommonErrServerErr
	}

	hash := common.HexToHash(txHash)

	tx, isPending, err := client.TransactionByHash(nil, hash)
	if err != nil {
		fmt.Println("查询失败:", err)
		return statecode.CommonErrServerErr
	}

	receipt, err := client.TransactionReceipt(nil, hash)
	if err != nil {
		fmt.Println("获取回执失败:", err)
		return statecode.CommonErrServerErr
	}

	*result = append(*result, models.Transaction{

		Hash:             tx.Hash().String(),
		BlockHash:        receipt.BlockHash.String(),
		BlockNumber:      receipt.BlockNumber.Uint64(),
		TransactionIndex: receipt.TransactionIndex,
		To:               tx.To().String(),
		Nonce:            tx.Nonce(),
		GasUsed:          receipt.GasUsed,
		Type:             tx.Type(),
	})

	db.Mysql.Table("txdatas").Create(&result)
	return statecode.CommonSuccess
}
