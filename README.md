# Raft Notes

## Introduction

Raft separates leader election, log replication and safety.
Raft states space reduction.

## Replicated State Machine

Consensus algorithms should have:
- safety: never return an incorrect result 安全
- functional: available 可用
- not depend on timing to ensure consistency 不依赖时钟
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


## Log Replication

Once a leader has been elected, it begins servicing client requests. Each client request contains a command to be executed by the replicated state machines. The leader appends the command to its log as a new entry, then issues AppendEntries RPCs in parallel to each of the other servers to replicate the entry. 
When the entry has been safely replicated (as described below), the leader applies the entry to its state machine and returns the result of that execution to the client. If followers crash or run slowly, or if network packets are lost, the leader retries AppendEntries RPCs indefinitely (even after it has responded to the client) until all followers eventually store all log entries.

产生leader后，leader开始处理 client 请求。每个请求包含一个 command，这个cmd将会被所有follower执行。 leader把log保存为新的entry, 然后并发 AppendEntries 到所有followers。当entry被安全复制到各个节点， leader把entry应用到自身节点，并把执行结果返回给客户端。如果 follower 崩溃或慢回应，leader无限重试， 直到所有 follower 保存了所有entries.

The leader decides when it is safe to apply a log entry to the state machines; such an entry is called committed. 
Raft guarantees that committed entries are durable and will eventually be executed by all of the available state machines. 

Entry 包含一个 Term 数值，和一个 Command 命令。entry被commit表示， entry cmd被应用到 SM 是安全的。 
Raft能保证 commited entry 最终会被所有 SM 执行。当leader 创建了一个 entry， 并且被复制到大部分 followers之上之后，entry就被认为是 committed. 这也会提交leader之前的entry，包括之前的leader创建的entry。
细节的情况稍后讨论。

The leader keeps track of the highest index it knows to be committed, and it includes that index in future AppendEntries RPCs (including heartbeats) so that the other servers eventually find out. 
leader 维护他提交的最新index(Log Index?), 并且把index包含在 AppendEntrie 中，因此其他follower可以发现。

Once a follower learns that a log entry is committed, it applies the entry to its local state machine (in log order).
当follower发现一个entry被提交了，他会应用这个entry到本地（垵序）。


Raft maintains the following properties, which together constitute the Log Matching Property in Figure 3:
• If two entries in different logs have the same index and term, then they store the same command.
不同log中的 index, term 都相同的 两个 entry，他们的命令相同
• If two entries in different logs have the same index and term, then the logs are identical in all preceding entries.
不同log中的 index, term 都相同的 两个 entry，他们之前的所有entry都相同

When sending an AppendEntries RPC, the leader includes the index and term of the entry in its log that immediately precedes the new entries. If the follower does not find an entry in its log with the same index and term, then it refuses the new entries. 
当leader发送 AppendEntries 时，包含了前一个entry的index和term. 如果follower 在自己的log中没有找到 这个 term, index的entry，则拒绝新的entry。这能保证第二个特点。

In Raft, the leader handles inconsistencies by forcing the followers’ logs to duplicate its own. This means that conflicting entries in follower logs will be overwritten with entries from the leader’s log. Section 5.4 will show that this is safe when coupled with one more restriction
follower可能缺少entry，多出entry，或者 both。 Raft处理不一致性的方法是，强制 follwer复制(overwrite) leader的log。

To bring a follower’s log into consistency with its own, the leader must find the latest log entry where the two logs agree, delete any entries in the follower’s log after that point, and send the follower all of the leader’s entries after that point. 
要保证follower的log的一致性，leader必须找到他们之前最近的那一条一致的log，删除follower的这一条log之后的log， 再发送leader的这一log之后的log给follower。这些动作是 AppendEntries 检查一致性（发现不一致）之后发生的。
leader为每个follower维护一个 nextIndex, 表示下一个准备发送给follower的 log entry。 当节点刚成为leader 的时候，他初始化 nextIndex 为 他持有的最新的log加一。 如果follower的log不一致，那么 AppendEntries的一致性检查在接下来的一次RPC中会失败。 PRC拒绝后，leader把nextIndex 减一，然后重试， 直到成功为止。这时，把follower上不一致的log删除，把leader的log追加到follower上。 这样 follower 就和 leader 一直了。Leader不覆盖或删除自己的log。

> If desired, the protocol can be optimized to reduce the number of rejected AppendEntries RPCs. For example, when rejecting an AppendEntries request, the follower can include the term of the conflicting entry and the first index it stores for that term. With this information, the leader can decrement nextIndex to bypass all of the conflicting entries in that term; one AppendEntries RPC will be required for each term with conflicting entries, rather than one RPC per entry. In practice, we doubt this optimization is necessary, since failures happen infrequently and it is unlikely that there will be many inconsistent entries.


## Safety

以上描述的机制还不够充分保证安全性。例如，follower错过了leader提交的几个entry，但是之后他立刻被选举为新的leader，这时他会覆盖掉其他server的log，导致其他server执行不同的命令顺序。
被选举leader的过程，需要引入一个限制。这个限制保证在任何一个term的leader拥有之前的term提交的所有entry。 


### Election restriction

Raft用简单的方式保证 在选举时，所有之前term的提交的log都存在于新的leader。log entry只能从 leader 单向流向 follower，leader不覆盖，删除他自己的log。
Raft使用投票过程来防止一个 candidate 赢得选举，除非他的log持有所有提交的entry。一个 candidate 必须解除大部分节点来赢得选举，这意味着每个提交的entry一定至少在其中一个server。 

If the candidate’s log is at least as up-to-date as any other log in that majority (where “up-to-date” is defined precisely below), then it will hold all the committed entries. The RequestVote RPC implements this restriction: the RPC includes information about the candidate’s log, and the voter denies its vote if its own log is more up-to-date than that of the candidate.
如果candidate的log至少和其他多数节点log一样"新" (up-to-date), 那么他将拥有所有提交的entry。 RequestVote 实现了这个约束： RPC包含了candidate的log信息，如果投票者发现自己的log比candidate更“新”， 那么投票者会拒绝给他投票。
“新”的定义是： 比较最新的entry的 index, term. term比较大的那个log 比较新， 如果term相同，那么 index 比较大的log 比较新。


### Committing entries from previous terms

As described in Section 5.3, a leader knows that an entry from its current term is committed once that entry is stored on a majority of the servers. 
If a leader crashes before committing an entry, future leaders will attempt to finish replicating the entry. 
However, a leader cannot immediately conclude that an entry from a previous term is committed once it is stored on a majority of servers. Figure 8 illustrates a situation where an old log entry is stored on a majority of servers, yet can still be overwritten by a future leader.

leader认为一个entry被复制到大多数节点后，就认为这个entry被提交了。如果leader在提交entry之前崩溃了，未来的 leader会尝试复制entry。但是leader无法立即得到结论，一个上一个term的entry是否被提交了. 
所以， Raft 不从过去的term中的log entry记录复制份数 并 提交。只对当前的term中entry 记录复制数，并提交entry。 当一个当前term的entry被提交了，那么他之前的所有entry都会被间接提交。

Raft 在提交规则中引入这个额外的复杂度，是因为 当leader复制之前term的 entries， entries 能 保持他们原本的term number。 这样维护term number容易复查问题，同时也减少了 新leader发送的entry的量。

### Safety argument

Given the complete Raft algorithm, we can now argue more precisely that the Leader Completeness Property holds (this argument is based on the safety proof; see Section 9.2). We assume that the Leader Completeness Property does not hold, then we prove a contradiction. 

Suppose the leader for term T (leaderT) commits a log entry from its term, but that log entry is not stored by the leader of some future term. Consider the smallest term U > T whose leader (leaderU) does not store the entry.

### Follower and candidate crashes

如果follower或者candidate奔溃，RequestVote and AppendEntries 会失败，Raft 会无限重试.
如果节点在响应RPC之前挂了，那么他重启之后会再次收到相同的RPC，因为 Raft 是幂等的， 所以重复RPC没有问题。 如果follower接收到 AppendEntries，发现他的log entry已经存在于自己的log中，那么直接忽略这个entry就行。

### Timing and availability

One of our requirements for Raft is that safety must not depend on timing: the system must not produce incorrect results just because some event happens more quickly or slowly than expected. 
However, availability (the ability of the system to respond to clients in a timely manner) must inevitably depend on timing. For example, if message exchanges take longer than the typical time between server crashes, candidates will not stay up long enough to win an election; without a steady leader, Raft cannot make progress.
Leader election is the aspect of Raft where timing is most critical. Raft will be able to elect and maintain a steady leader as long as the system satisfies the following timing requirement:

> broadcastTime << electionTimeout << MTBF

In this inequality broadcastTime is the average time it takes a server to send RPCs in parallel to every server in the cluster and receive their responses; 
broadcastTime 指 一个server发送RPC的平均间隔时间。

electionTimeout is the election timeout described in Section 5.2; 
停留在follower 状态的时间，150 - 300 ms 随机

and MTBF is the average time between failures for a single server. 
单个节点的两次崩溃之间的运行时间。

broadcast time 应该比 electionTimeout 小一个数量级。
当leader奔溃，系统会有大约一个 electionTimeout 时长的无法服务。
broadcast time may range from 0.5ms to 20ms, election timeout is likely to be somewhere between 10ms and 500ms.


## Cluster membership changes

to replace servers when they fail or to change the degree of replication，we decided to automate configuration changes and incorporate them into the Raft consensus algorithm.
有时我们需要替换节点，或改变节点数。Raft有自动化配置变更的机制。

不可能一次自动改变所有节点的配置，所以，集群可以在转变时分成两个独立的部分。为了保证安全，配置变更必须分为2个步骤。 集群首先转变使用一个 变化配置， 叫做 joint consensus; 第二步，系统再转变为新的配置。 joint consensus 组合了 旧配置 和 新配置。 

1. Log entries are replicated to all servers in both configurations. 
log entry 被复制到了新旧配置中的所有节点。

2. Any server from either configuration may serve as leader. 
新旧配置中的任意节点可能成为leader

3. Agreement (for elections and entry commitment) requires separate majorities from both the old and new configurations.  
选举和提交日志需要 旧配置中的大多数节点 和 新配置中的大多数节点 的同意。








