package models

import (
	"gorm.io/gorm"
	"math/big"
	"pledge-backend/db"
	"time"
)

type Transaction struct {
	// 交易基本信息
	Hash                 string    `json:"hash"`                           // 交易哈希，唯一标识
	BlockHash            string    `json:"blockHash"`                      // 包含该交易的区块哈希
	BlockNumber          uint64    `json:"blockNumber"`                    // 区块高度
	TransactionIndex     uint      `json:"transactionIndex"`               // 交易在区块中的索引位置
	Timestamp            time.Time `json:"timestamp"`                      // 交易时间戳
	From                 string    `json:"from"`                           // 发送方地址
	To                   string    `json:"to,omitempty"`                   // 接收方地址（合约创建交易可能为空）
	Value                *big.Int  `json:"value"`                          // 交易金额（以最小单位计，如Wei）
	Nonce                uint64    `json:"nonce"`                          // 发送方账户交易计数器
	GasLimit             uint64    `json:"gasLimit"`                       // 交易允许消耗的最大Gas
	GasPrice             *big.Int  `json:"gasPrice"`                       // Gas价格
	GasUsed              uint64    `json:"gasUsed,omitempty"`              // 实际消耗的Gas（执行后确定）
	Input                string    `json:"input"`                          // 交易数据（如合约调用参数）
	Status               uint      `json:"status,omitempty"`               // 交易状态（0失败，1成功）
	Type                 uint8     `json:"type"`                           // 交易类型（如Legacy, EIP-1559等）
	MaxPriorityFeePerGas *big.Int  `json:"maxPriorityFeePerGas,omitempty"` // 最大优先费（EIP-1559）
	MaxFeePerGas         *big.Int  `json:"maxFeePerGas,omitempty"`         // 最大总费用（EIP-1559）
}

type TxData struct {
	gorm.Model
	Hash                 string `gorm:"column:hash;type:char(66);primaryKey;not null" json:"hash"`
	BlockHash            string `gorm:"column:block_hash;type:char(66);index" json:"blockHash"`
	BlockNumber          uint64 `gorm:"column:block_number;index" json:"blockNumber"`
	TransactionIndex     uint   `gorm:"column:transaction_index" json:"transactionIndex"`
	Timestamp            int64  `gorm:"column:timestamp;index" json:"timestamp"`
	From                 string `gorm:"column:from_address;type:char(42);index" json:"from"`
	To                   string `gorm:"column:to_address;type:char(42);index" json:"to,omitempty"`
	Value                string `gorm:"column:value;type:varchar(78);not null" json:"value"`
	Nonce                uint64 `gorm:"column:nonce" json:"nonce"`
	GasLimit             uint64 `gorm:"column:gas_limit" json:"gasLimit"`
	GasPrice             string `gorm:"column:gas_price;type:varchar(78);not null" json:"gasPrice"`
	GasUsed              uint64 `gorm:"column:gas_used" json:"gasUsed,omitempty"`
	Input                string `gorm:"column:input_data;type:text" json:"input"`
	Status               uint   `gorm:"column:status" json:"status,omitempty"`
	Type                 uint8  `gorm:"column:type" json:"type"`
	MaxPriorityFeePerGas string `gorm:"column:max_priority_fee_per_gas;type:varchar(78)" json:"maxPriorityFeePerGas,omitempty"`
	MaxFeePerGas         string `gorm:"column:max_fee_per_gas;type:varchar(78)" json:"maxFeePerGas,omitempty"`
}

func NewTxData() *TxData {
	return &TxData{}
}

func (b *TxData) TableName() string {
	return "txdatas"
}

func (b *TxData) TxData(txHash string, res *[]Transaction) error {

	var txData []TxData

	err := db.Mysql.Table("txdatas").Where("hash=?", txHash).Find(&txData).Debug().Error

	if err != nil {
		return err
	}

	for _, v := range txData {
		*res = append(*res, Transaction{
			Hash:             v.Hash,
			BlockHash:        v.BlockHash,
			BlockNumber:      v.BlockNumber,
			TransactionIndex: v.TransactionIndex,
			From:             v.From,
			To:               v.To,
			Nonce:            v.Nonce,
			GasLimit:         v.GasLimit,
			GasUsed:          v.GasUsed,
			Input:            v.Input,
			Status:           v.Status,
			Type:             v.Type,
		})
	}

	return nil
}
