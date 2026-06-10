package raftadapter

// PeerFactory describes the minimal peer construction contract required by the
// business layer. It keeps the application code independent from the Raft
// implementation package and only exposes the stable API required by the
// design.
type PeerFactory interface {
	GetState() (term int, isLeader bool)
	Start(command interface{}) (index int, term int, isLeader bool)
	Kill()
}

// ApplyObserver represents the minimal callback surface for applying committed
// commands from Raft into the business layer.
type ApplyObserver interface {
	HandleApply(msg ApplyMsg)
}
