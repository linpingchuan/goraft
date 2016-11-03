# Raft Notes

## Introduction

Raft separates leader election, log replication and safety.
Raft states space reduction.

## Replicated State Machine

Consensus algorithms should have:
- safety: never return an incorrect result 安全
- functional: available 可用
- not depend on timing to unsure consistency 不依赖时钟
- a minority of slow server should not impact the whole system 少数慢节点不影响log命令被执行

## Design 

Distinguished Leader: Complete responsibility form managing the replicated log.
Leader负责管理log

accepts log entries from clients, replicates them on other servers, and tells servers when it is safe to apply log entries to their state machines
接受其他节点的log，复制到其他节点，并告知他们何时可以安全执行

- leader election : a new leader must be chosen when an existing leader fails
失败时选举新的Leader

- log replication : the leader must accept log entries from clients and replicate them across the cluster, forcing the other logs to agree with its own
接受Log，复制Log

- safety : if any server has applied a particular log entry to its state machine, then no other server may apply a different command for the same log index
如果任意一个server执行了某个log, 那么对于同一个 log index, 其他的节点不能执行不同的log command. 
(例如 节点A 执行了 {id:13, cmd: "remove 3"}, 那么其他节点不能执行 id=13的其他命令)


### Safety Property:
关于安全性的概述

- Election Safety: at most one leader can be elected in a given term.
最多一个Leader

- Leader Append-Only: a leader never overwrites or deletes entries in its log; it only appends new entries.
Leader只能 append log, 不能 update 或 delete

- Log Matching: if two logs contain an entry with the same index and term, then the logs are identical in all entries up through the given index.
如果 Log 拥有 相同的 index, term, 则这两个log是完全相同的

- Leader Completeness: if a log entry is committed in a given term, then that entry will be present in the logs of the leaders for all higher-numbered terms.
在某一轮(term)选举后，如果log被 commit, 那么在以后的更高轮次选举后，这个log还是存在的

- State Machine Safety: if a server has applied a log entry at a given index to its state machine, no other server will ever apply a different log entry for the same index.
如果任意一个server执行了某个log, 那么对于同一个 log index, 其他的节点不能执行不同的log command. 


## Basics

Server has 3 States: Leader, Follower, Candidate

Transition:

Start up -> Follower 
Follower times out, start election -> Candicate
Candidate time out, start new election -> Candidate
Candidate receives votes from majority -> Leader
Candidate discover Leader or new term -> Follower
Leader discovers higher term -> Follower

Terms:
begins with election,consecutive integer 选举产生轮次, 是连续整数
if split vote happens, new election begins shortly. 如果平票，发起新选举



