package raftadapter

import "context"

type ApplyLoop struct {
	applyCh    <-chan ApplyMsg
	dispatcher ApplyObserver
}

func NewApplyLoop(applyCh <-chan ApplyMsg, dispatcher ApplyObserver) *ApplyLoop {
	return &ApplyLoop{applyCh: applyCh, dispatcher: dispatcher}
}

func (l *ApplyLoop) Run(ctx context.Context) error {
	if l == nil || l.dispatcher == nil {
		return nil
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-l.applyCh:
			if !ok {
				return nil
			}
			l.dispatcher.HandleApply(msg)
		}
	}
}
