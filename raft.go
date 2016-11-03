package goraft

import "sync/atomic"

const (
	FollowState uint32 = iota
	CandidateState
	LeaderState
)

type Raft struct {
	State    uint32
	CurrTerm uint32

	Host  string
	Peers []string
}

func (r *Raft) updateState(state uint32) {
	atomic.StoreUint32(&r.State, state)
}

func (r *Raft) getState() uint32 {
	return atomic.LoadUint32(&r.State)
}
