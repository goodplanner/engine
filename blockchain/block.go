package blockchain

import (
	"github.com/it-chain/yggdrasill"
	"github.com/it-chain/yggdrasill/common"
	"github.com/it-chain/yggdrasill/impl"
	"github.com/it-chain/midgard"
	"github.com/it-chain/it-chain-Engine/core/eventstore"
)

type Block = common.Block

type DefaultBlock = impl.DefaultBlock

type Repository interface {
	yggdrasill.BlockStorageManager
	NewEmptyBlock() (Block, error)
	GetBlockCreator() string
}

type BlockHeight = uint64

type BlockPool interface {
	Add(block Block)
	Get(height BlockHeight) Block
	Delete(height BlockHeight)
}


// block queued Aggregate id
var BLOCK_QUEUED_AID = "BLOCK_QUEUED_AID"

type BlockPoolModel struct {
	pool map[BlockHeight] Block
}

func (p *BlockPoolModel) Add(block Block) {
	height := block.GetHeight()
	p.pool[height] = block

	event := createBlockQueuedEvent(block)

	eventstore.Save(BLOCK_QUEUED_AID, event)
}

func (p *BlockPoolModel) Get(height BlockHeight) Block {
	return p.pool[height]
}

func (p *BlockPoolModel) Delete(height BlockHeight) {
	delete(p.pool, height)
}

func createBlockQueuedEvent(block Block) BlockQueuedEvent {
	return BlockQueuedEvent{
		EventModel: midgard.EventModel{
			ID: BLOCK_QUEUED_EID,
		},
		Block: block,
	}
}


// BlockSyncState Aggregate ID
var BC_SYNC_STATE_AID = "BC_SYNC_STATE_AID"

// 현재 블록 동기화가 진행 중인지 정보를 가진다.
type BlockSyncState struct {
	midgard.AggregateModel
	isProgress bool
}

func (state BlockSyncState) GetID() string {
	return BC_SYNC_STATE_AID
}

func (state BlockSyncState) On(event midgard.Event) error {
	// TODO: Event에 따라서 isProgress 값 바꿔주기
	return nil
}

