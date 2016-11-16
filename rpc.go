package goraft

// type AERequest struct {
// 	LeaderID     string
// 	Term         uint32
// 	PrevLogTerm  uint32
// 	PrevLogIndex uint64
// 	LeaderCommit uint64
// 	Entries      []LogEntry
// }

// type AEResp struct {
// 	Term    uint32
// 	Success bool
// }

// type RVRequest struct {
// 	CandidateID  string
// 	LastLogIndex uint64
// 	LastLogTerm  uint32
// 	Term         uint32
// }

// type RVResp struct {
// 	Term        uint32
// 	VoteGranted bool
// }

// AppendEntries to followers
func (raft *Raft) AppendEntries(AERequest) AEResp {

	return AEResp{}
}

// RequestVote for leader election
func (raft *Raft) RequestVote(RVRequest) RVResp {

	return RVResp{}
}
