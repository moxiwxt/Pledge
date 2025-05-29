package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"pledge-backend/db"
	"pledge-backend/log"
	"time"
)

func NewBlockData() *BlockInfo {
	return &BlockInfo{}
}

// BlockInfo 区块信息结构
type BlockInfo struct {
	Number     uint64 `json:"number"`
	Hash       string `json:"hash"`
	Timestamp  uint64 `json:"timestamp"`
	Difficulty string `json:"difficulty"`
}

// BlockType 区块类型枚举
type BlockType int

const (
	HeadBlock BlockType = iota
	FinalizedBlock
	SafeBlock
)

// fetchAndSaveBlocks 获取并保存区块信息
func (s *BlockInfo) FetchAndSaveBlocks() {
	// 创建区块信息映射
	blockMap := make(map[BlockType]*types.Block)

	ethClient, err := ethclient.Dial("xxxxxxxxxxxxx")
	ctx := context.Background()

	// 获取不同类型的区块
	header, err := ethClient.HeaderByNumber(ctx, nil)
	if err != nil {
		log.Logger.Error(err.Error())
	}

	headBlock, err := ethClient.BlockByHash(ctx, header.Hash())
	if err != nil {
		log.Logger.Error(err.Error())
	}
	blockMap[HeadBlock] = headBlock

	finalizedBlock, err := getFinalizedBlock(ctx, ethClient)
	if err != nil {
		log.Logger.Error(err.Error())
	}
	blockMap[FinalizedBlock] = finalizedBlock

	safeBlock, err := ethClient.BlockByNumber(ctx, nil)
	if err != nil {
		log.Logger.Error(err.Error())
	}
	blockMap[SafeBlock] = safeBlock

	// 使用switch case处理不同类型的区块
	for blockType, block := range blockMap {
		var keyPrefix string

		switch blockType {
		case HeadBlock:
			keyPrefix = "head_block"
		case FinalizedBlock:
			keyPrefix = "finalized_block"
		case SafeBlock:
			keyPrefix = "safe_block"
		default:
			log.Logger.Info("未知区块类型:")
			continue
		}

		// 保存区块信息到Redis
		if err := saveBlockToRedis(ctx, keyPrefix, block); err != nil {
			log.Logger.Error(err.Error())
		} else {
			log.Logger.Info("%s区块信息已更新")
		}
	}

}

// getFinalizedBlock 获取finalized区块
func getFinalizedBlock(ctx context.Context, ethClient *ethclient.Client) (*types.Block, error) {
	// 实际应用中应使用客户端特定的API获取finalized区块
	// 这里使用简化逻辑，假设finalized是当前区块的前5个区块
	header, err := ethClient.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}

	finalizedNumber := big.NewInt(0).Sub(header.Number, big.NewInt(5))
	if finalizedNumber.Cmp(big.NewInt(0)) < 0 {
		finalizedNumber = big.NewInt(0)
	}

	return ethClient.BlockByNumber(ctx, finalizedNumber)
}

// saveBlockToRedis 保存区块信息到Redis
func saveBlockToRedis(ctx context.Context, keyPrefix string, block *types.Block) error {
	blockInfo := BlockInfo{
		Number:     block.NumberU64(),
		Hash:       block.Hash().String(),
		Timestamp:  block.Time(),
		Difficulty: block.Difficulty().String(),
	}

	// 序列化为JSON
	jsonData, err := json.Marshal(blockInfo)
	if err != nil {
		return fmt.Errorf("序列化区块信息失败: %v", err)
	}

	var result map[string]string
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return fmt.Errorf("区块信息转map失败: %v", err)
	}

	// 使用Redis Hash存储区块信息
	err1 := db.RedisSetHash(keyPrefix, result, 24*time.Hour)
	if err1 != nil {
		return err1
	}

	return nil
}
