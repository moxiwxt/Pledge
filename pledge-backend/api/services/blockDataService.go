package services

import (
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"pledge-backend/api/common/statecode"
	"pledge-backend/api/models"
	"pledge-backend/db"
	"pledge-backend/log"
)

type BlockDataService struct{}

func NewGetBlockData() *BlockDataService { return &BlockDataService{} }

func (s *BlockDataService) BlockData(blockNum string, full string, result *[]models.BlockDataRes) int {

	// 检查是否为特殊区块标识
	isSpecial := blockNum == "head" || blockNum == "finalized" || blockNum == "safe"

	if isSpecial {
		// 特殊区块标识，直接查询最新的区块
		bool := db.RedisExists(blockNum)

		if bool {
			data, err := db.RedisGet(blockNum)
			if err != nil {
				log.Logger.Error(err.Error())
				return statecode.CommonErrServerErr
			}
			_ = json.Unmarshal(data, result)
		} else {
			client, err := ethclient.Dial("")
			if err != nil {
				log.Logger.Error(err.Error())
				return statecode.CommonErrServerErr
			}
			var num *big.Int
			_, success := num.SetString(blockNum, 10)
			if success {
				//header, err := client.HeaderByNumber(context.Background(), num)
				//if err != nil {
				//	log.Logger.Error(err.Error())
				//	return statecode.CommonErrServerErr
				//}
				block, err := client.BlockByNumber(context.Background(), num)
				if err != nil {
					log.Logger.Error(err.Error())
					return statecode.CommonErrServerErr
				}

				for _, tx := range block.Transactions() {
					*result = append(*result, models.BlockDataRes{
						BlockData: models.BlockResult{
							BlobVersionedHashes:  tx.BlobVersionedHashes,
							BlockHash:            tx.BlockHash,
							BlockNumber:          tx.BlockNumber,
							ChainID:              tx.ChainID,
							From:                 tx.From,
							Gas:                  tx.Gas,
							GasPrice:             tx.GasPrice,
							Hash:                 tx.Hash,
							Input:                tx.Input,
							MaxFeePerBlobGas:     tx.MaxFeePerBlobGas,
							MaxFeePerGas:         v.MaxFeePerGas,
							MaxPriorityFeePerGas: v.MaxPriorityFeePerGas,
							Nonce:                tx.Nonce,
							R:                    tx.R,
							S:                    tx.S,
							To:                   tx.To,
							TransactionIndex:     tx.Hash,
							Type:                 tx.Type,
							V:                    tx.V,
							Value:                tx.Value,
							YParity:              tx.YParity,
						},
					})
				}

				db.Mysql.Table("blockdatas").Create(&result)
			}
			return statecode.CommonSuccess
		}
	} else {
		// 普通区块号查询
		err := models.NewBlockData().BlockResult(blockNum, result)
		if err != nil {
			log.Logger.Error(err.Error())
			return statecode.CommonErrServerErr
		}
		return statecode.CommonSuccess
	}

}
