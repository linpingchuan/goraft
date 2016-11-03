package goraft

func startElectionTimeoutTicker() {
	// when is follower
	// set a ticker, reset it when receive heartbeat
	// if timeout , set state to candidate
}

func stopElectionTimeoutTicker() {
	// when not follower

}

func becomeFollower() {
	// set state

	// startElectionTimeoutTicker
}
