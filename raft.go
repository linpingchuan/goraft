package goraft

const (
	FollowState uint32 = iota
	CandidateState
	LeaderState
)

// Raft state machine
type Raft struct {
	ID string
	// State
	State int
	// Peers hosts
	Peers []string

	// persistent on all server
	CurrTerm uint32
	VotedFor string
	Logs     []LogEntry

	// volatile on all server
	commitIndex uint64
	lastApplied uint64

	// volatile on Leader
	nextIndex  []uint64
	matchIndex []uint64
}
