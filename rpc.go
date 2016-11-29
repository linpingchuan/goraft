package goraft

import context "golang.org/x/net/context"
import "log"

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
func (raft *Raft) AppendEntries(ctx context.Context, req *AERequest) (*AEResp, error) {
	log.Printf("%#v", *req)

	return &AEResp{}, nil
}

// RequestVote for leader election
func (raft *Raft) RequestVote(ctx context.Context, req *RVRequest) (*RVResp, error) {
	log.Printf("%#v", *req)

	return &RVResp{}, nil
}
