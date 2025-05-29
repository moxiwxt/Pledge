package models

import (
	"pledge-backend/db"
)

// 根结构体
type BlockResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  BlockResult `json:"result"`
	ID      int         `json:"id"`
}

// 区块结果结构体
type BlockResult struct {
	BlockID              int      `json:"block_id"`
	AccessList           []string `json:"accessList"`
	BlobVersionedHashes  []string `json:"blobVersionedHashes"`
	BlockHash            string   `json:"blockHash"`
	BlockNumber          string   `json:"blockNumber"`
	ChainID              string   `json:"chainId"`
	From                 string   `json:"from"`
	Gas                  string   `json:"gas"`
	GasPrice             string   `json:"gasPrice"`
	Hash                 string   `json:"hash"`
	Input                string   `json:"input"`
	MaxFeePerBlobGas     string   `json:"maxFeePerBlobGas"`
	MaxFeePerGas         string   `json:"maxFeePerGas"`
	MaxPriorityFeePerGas string   `json:"maxPriorityFeePerGas"`
	Nonce                string   `json:"nonce"`
	R                    string   `json:"r"`
	S                    string   `json:"s"`
	To                   string   `json:"to"`
	TransactionIndex     string   `json:"transactionIndex"`
	Type                 string   `json:"type"`
	V                    string   `json:"v"`
	Value                string   `json:"value"`
	YParity              string   `json:"yParity"`
}

// BlockResult GORM模型（对应数据库表）
type BlockData struct {
	Id                   int    `json:"-" gorm:"column:id;primaryKey;autoIncrement"`
	BlockID              int    `json:"block_id" gorm:"column:block_id;"`
	BlockHash            string `gorm:"column:block_hash;type:varchar(66);index;not null" json:"blockHash"`
	BlockNumber          string `gorm:"column:block_number;type:varchar(20);index;not null" json:"blockNumber"`
	ChainID              string `gorm:"column:chain_id;type:varchar(20);not null" json:"chainId"`
	From                 string `gorm:"column:from_address;type:varchar(42);index;not null" json:"from"`
	Gas                  string `gorm:"column:gas;type:varchar(20);not null" json:"gas"`
	GasPrice             string `gorm:"column:gas_price;type:varchar(20);not null" json:"gasPrice"`
	Hash                 string `gorm:"column:tx_hash;type:varchar(66);uniqueIndex;not null" json:"hash"`
	Input                string `gorm:"column:input_data;type:text;not null" json:"input"`
	MaxFeePerBlobGas     string `gorm:"column:max_fee_per_blob_gas;type:varchar(20);not null" json:"maxFeePerBlobGas"`
	MaxFeePerGas         string `gorm:"column:max_fee_per_gas;type:varchar(20);not null" json:"maxFeePerGas"`
	MaxPriorityFeePerGas string `gorm:"column:max_priority_fee_per_gas;type:varchar(20);not null" json:"maxPriorityFeePerGas"`
	Nonce                string `gorm:"column:nonce;type:varchar(20);not null" json:"nonce"`
	R                    string `gorm:"column:signature_r;type:varchar(132);not null" json:"r"`
	S                    string `gorm:"column:signature_s;type:varchar(132);not null" json:"s"`
	To                   string `gorm:"column:to_address;type:varchar(42);index;not null" json:"to"`
	TransactionIndex     string `gorm:"column:transaction_index;type:varchar(20);not null" json:"transactionIndex"`
	Type                 string `gorm:"column:tx_type;type:varchar(10);not null" json:"type"`
	V                    string `gorm:"column:signature_v;type:varchar(10);not null" json:"v"`
	Value                string `gorm:"column:value;type:varchar(20);not null" json:"value"`
	YParity              string `gorm:"column:y_parity;type:varchar(10);not null" json:"yParity"`
}

// 提现信息结构体
type Withdrawal struct {
	Address        string `json:"address"`
	Amount         string `json:"amount"`
	Index          string `json:"index"`
	ValidatorIndex string `json:"validatorIndex"`
}

type BlockDataRes struct {
	Index     int         `json:"index"`
	BlockData BlockResult `json:"block_data"`
}

func NewBlockData() *BlockData {
	return &BlockData{}
}

func (b *BlockData) TableName() string {
	return "blockdatas"
}

func (b *BlockData) BlockResult(blockNum string, res *[]BlockDataRes) error {

	var blockData []BlockData

	err := db.Mysql.Table("blockdatas").Where("block_number=?", blockNum).Order("chain_id asc").Find(&blockData).Debug().Error

	if err != nil {
		return err
	}

	for _, v := range blockData {
		*res = append(*res, BlockDataRes{
			Index: v.BlockID - 1,
			BlockData: BlockResult{
				BlockID:              v.BlockID,
				BlockHash:            v.BlockHash,
				BlockNumber:          v.BlockNumber,
				ChainID:              v.ChainID,
				From:                 v.From,
				Gas:                  v.Gas,
				GasPrice:             v.GasPrice,
				Hash:                 v.Hash,
				Input:                v.Input,
				MaxFeePerBlobGas:     v.MaxFeePerBlobGas,
				MaxFeePerGas:         v.MaxFeePerGas,
				MaxPriorityFeePerGas: v.MaxPriorityFeePerGas,
				Nonce:                v.Nonce,
				R:                    v.R,
				S:                    v.S,
				To:                   v.To,
				TransactionIndex:     v.TransactionIndex,
				Type:                 v.Type,
				V:                    v.V,
				Value:                v.Value,
				YParity:              v.YParity,
			},
		})
	}
	return nil

}
