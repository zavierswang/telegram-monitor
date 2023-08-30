package jobs

import "context"

type ListenTRC10 struct {
	ticker     int
	trc10Queue []string
	ctx        context.Context
}

func (l *ListenTRC10) Run() {

}
