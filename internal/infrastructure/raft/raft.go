package raftadapter

import (
	"raft-stash/labrpc"
	raft "raft-stash/raft"
)

type ApplyMsg = raft.ApplyMsg

type Persister = raft.Persister

type ClientEnd = labrpc.ClientEnd

type Raft interface {
	GetState() (int, bool)
	Start(command any) (int, int, bool)
	Kill()
}

type Adapter struct {
	inner *raft.Raft
}

func New(peers []*labrpc.ClientEnd, me int, persister *raft.Persister, applyCh chan raft.ApplyMsg) *Adapter {
	return &Adapter{inner: raft.Make(peers, me, persister, applyCh)}
}

func (a *Adapter) GetState() (int, bool) {
	return a.inner.GetState()
}

func (a *Adapter) GetLeader() string {
	_, isLeader := a.GetState()
	if !isLeader {
		return ""
	}
	return "self"
}

func (a *Adapter) Start(command any) (int, int, bool) {
	return a.inner.Start(command)
}

func (a *Adapter) Kill() {
	a.inner.Kill()
}

func (a *Adapter) Inner() *raft.Raft {
	return a.inner
}
