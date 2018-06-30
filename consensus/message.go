package consensus

import (
	"github.com/it-chain/it-chain-Engine/blockchain"
	"encoding/json"
	"sync"
)

type PrePrepareMsg struct {
	ConsensusId   ConsensusId
	SenderId      string
	ProposedBlock blockchain.Block
}

func (pp PrePrepareMsg) ToByte() ([]byte, error) {
	data, err := json.Marshal(pp)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type PrepareMsg struct {
	ConsensusId   ConsensusId
	SenderId      string
	ProposedBlock blockchain.Block
}

func (p PrepareMsg) ToByte() ([]byte, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type CommitMsg struct {
	ConsensusId ConsensusId
	SenderId    string
}

func (c CommitMsg) ToByte() ([]byte, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type PrepareMsgRepository interface {
	DeleteAllPrepareMsg(consensusId ConsensusId)
	InsertPrepareMsg(prepareMsg PrepareMsg)
	FindPrepareMsgsByConsensusID(consensusId ConsensusId) []PrepareMsg
}

type PrepareMsgRepositoryImpl struct {
	PrepareMsgPool map[ConsensusId][]PrepareMsg
	lock           *sync.RWMutex
}

type CommitMsgRepository interface {
	DeleteAllCommitMsg(consensusId ConsensusId)
	InsertCommitMsg(commitMsg CommitMsg)
	FindCommitMsgsByConsensusID(consensusId ConsensusId) []CommitMsg
}

type CommitMsgRepositoryImpl struct {
	CommitMsgPool map[ConsensusId][]CommitMsg
	lock          *sync.RWMutex
}

func NewPrepareMsgRepository() *PrepareMsgRepositoryImpl {
	return &PrepareMsgRepositoryImpl{
		PrepareMsgPool: make(map[ConsensusId][]PrepareMsg),
		lock:           &sync.RWMutex{},
	}
}

func NewCommitMsgRepository() *CommitMsgRepositoryImpl {
	return &CommitMsgRepositoryImpl{
		CommitMsgPool: make(map[ConsensusId][]CommitMsg),
		lock:          &sync.RWMutex{},
	}
}

func (pr *PrepareMsgRepositoryImpl) DeleteAllPrepareMsg(consensusId ConsensusId) {
	pr.lock.Lock()
	defer pr.lock.Unlock()

	delete(pr.PrepareMsgPool, consensusId)
}

func (cr *CommitMsgRepositoryImpl) DeleteAllCommitMsg(consensusId ConsensusId) {
	cr.lock.Lock()
	defer cr.lock.Unlock()

	delete(cr.CommitMsgPool, consensusId)
}

func (pr *PrepareMsgRepositoryImpl) InsertPrepareMsg(prepareMsg PrepareMsg) {
	pr.lock.Lock()
	defer pr.lock.Unlock()

	if prepareMsg.SenderId == "" {
		return
	}

	prepareMsgPool := pr.PrepareMsgPool[prepareMsg.ConsensusId]

	if prepareMsgPool == nil {
		pr.PrepareMsgPool[prepareMsg.ConsensusId] = make([]PrepareMsg, 0)
	}

	var hasItem = func(prepareMsgPool []PrepareMsg, prepareMsg PrepareMsg) bool {
		for _, msg := range prepareMsgPool {
			if msg.SenderId == prepareMsg.SenderId {
				return true
			}
		}
		return false
	}(prepareMsgPool, prepareMsg)

	if hasItem {
		return
	}

	pr.PrepareMsgPool[prepareMsg.ConsensusId] = append(pr.PrepareMsgPool[prepareMsg.ConsensusId], prepareMsg)
}

func (cr *CommitMsgRepositoryImpl) InsertCommitMsg(CommitMsg CommitMsg) {
	cr.lock.Lock()
	defer cr.lock.Unlock()

	if CommitMsg.SenderId == "" {
		return
	}

	CommitMsgPool := cr.CommitMsgPool[CommitMsg.ConsensusId]

	if CommitMsgPool == nil {
		cr.CommitMsgPool[CommitMsg.ConsensusId] = make([]CommitMsg, 0)
	}

	var hasItem = func(CommitMsgPool []CommitMsg, CommitMsg CommitMsg) bool {
		for _, msg := range CommitMsgPool {
			if msg.SenderId == CommitMsg.SenderId {
				return true
			}
		}
		return false
	}(CommitMsgPool, CommitMsg)

	if hasItem {
		return
	}

	cr.CommitMsgPool[CommitMsg.ConsensusId] = append(cr.CommitMsgPool[CommitMsg.ConsensusId], CommitMsg)
}

func (pr *PrepareMsgRepositoryImpl) FindPrepareMsgsByConsensusID(consensusId ConsensusId) []PrepareMsg {

	return pr.PrepareMsgPool[consensusId]
}

func (cr *CommitMsgRepositoryImpl) FindCommitMsgsByConsensusID(consensusId ConsensusId) []CommitMsg {

	return cr.CommitMsgPool[consensusId]
}