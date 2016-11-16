package goraft

type logEntry struct {
	Command []byte
	Index   uint64
	Term    uint32
}
