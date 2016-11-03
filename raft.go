package goraft

import "sync/atomic"

const (
	FollowState uint32 = iota
	CandidateState
	LeaderState
)

// Raft state machine
type Raft struct {
	// State
	State uint32
	// CurrTerm
	CurrTerm uint32

	// Host of current server
	Host string
	// Leader host
	Leader string
	// Peers hosts
	Peers []string
	// IsLeader
	IsLeader bool
}

func (r *Raft) updateState(state uint32) {
	atomic.StoreUint32(&r.State, state)
}

func (r *Raft) getState() uint32 {
	return atomic.LoadUint32(&r.State)
}
