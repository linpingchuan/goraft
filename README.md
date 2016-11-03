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

Terms act as a logical clock [14] in Raft, and they allow servers to detect obsolete information such as stale leaders.
Term 是一种 逻辑时钟, 它让server能够检测到过期的信息，例如 过期的Leader。 每个Server都会存一个当前term的数字，并且该值单调递增。

Current terms are exchanged whenever servers communicate; 
节点产生交互时，交换 current term

if one server’s current term is smaller than the other’s, then it updates its current term to the larger value. 
当节点发现自己的current term 小于 其他节点，则更新自己的值。

If a candidate or leader discovers that its term is out of date, it immediately reverts to fol- lower state. 
如果 candidate, leader 发现自己的 current term 小，则转换到 follower state.

If a server receives a request with a stale term number, it rejects the request.
如果server 接收到 小于自身 term 的请求，拒绝请求。


Two Basic RPCs:

- RequestVote:
initiated by candidates during elections 竞选

- AppendEntries:
initiated by leaders to replicate log entries and to provide a form of heartbeat
leader 用来复制日志 和 心跳

A third RPC is introduced latter for transferring snapshots between servers

Servers retry RPCs if they do not receive a response in a timely manner, and they issue RPCs in parallel for best performance.
如果RPC无响应，一定时间后重试；PRC并发。


## Leader Election

Raft uses a heartbeat mechanism to trigger leader election.

A server remains in follower state as long as it receives valid RPCs from a leader or candidate. Leaders send periodic heartbeats (AppendEntries RPCs that carry no log entries) to all followers in order to maintain their authority. If a follower receives no communication over a period of time called the election timeout, then it assumes there is no viable leader and begins an election to choose a new leader.
Leader周期性发送心跳给followers, 如果follower一段时间没有收到消息，发生 election timeout, 他会尝试新一轮选举。

To begin an election, a follower increments its current term and transitions to candidate state. It then votes for itself and issues RequestVote RPCs in parallel to each of the other servers in the cluster. A candidate continues in this state until one of three things happens: (a) it wins the election, (b) another server establishes itself as leader, or (c) a period of time goes by with no winner. These outcomes are discussed separately in the paragraphs below.
开始选举的操作，follower递增自身current term, 转换到 candidate state. 发起 RequestVote 给其他节点， 选举有三种结果：1 赢得选举； 2 别的节点赢得选举； 3 一段时间过去后还是没有节点成为leader

1. 赢得选举

A candidate wins an election if it receives votes from a majority of the servers in the full cluster for the same term. 
cadidate得到大多数成员的投票后，赢得选举。

Each server will vote for at most one candidate in a given term, on a first-come-first-served basis (note: Section 5.4 adds an additional restriction on votes). 
每一轮term, 只能投给一个candidate

The majority rule ensures that at most one candidate can win the election for a particular term (the Election Safety Prop- erty in Figure 3). 
每次最多一个节点赢得选举


2. 别的节点赢得选举

While waiting for votes, a candidate may receive an AppendEntries RPC from another server claiming to be leader. If the leader’s term (included in its RPC) is at least as large as the candidate’s current term, then the candidate recognizes the leader as legitimate and returns to follower state. If the term in the RPC is smaller than the candidate’s current term, then the candidate rejects the RPC and con- tinues in candidate state.
在等待vote结果时，candidate可能会收到其他节点的Leader心跳(AppendEntries), 这时需要比较 current term 和 AppendEntries 传来的 term, 如果 currTerm > term, 则 拒绝请求，否则，自身退回 follower state, 承认对方为 Leader.


3. 一段时间过去后还是没有节点成为leader

if many followers become candidates at the same time, votes could be split so that no candidate obtains a majority. When this happens, each candidate will time out and start a new election by incrementing its term and initiating another round of RequestVote RPCs. 
多个节点同时成为candidate, 那么可能没有节点获得大多数vote。那么等待下次 election timeout ，重新进行选举。为了防止不停重复这种情况， election timeout 的时长一般设定在 150 - 300 ms 随机。对于平票，也是一样，每个candidate在开始竞选时，随机一个 election timeout， 如果平票，等待timeout时长，再进行下一次选举，这样下次选举再次平票的可能性就比较小了。











