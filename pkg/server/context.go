package server

import "time"

type Context struct {
}

func (*Context) Deadline() (deadline time.Time, ok bool) {
	panic("implement me")
}

func (*Context) Done() <-chan struct{} {
	panic("implement me")
}

func (*Context) Err() error {
	return nil
}

func (*Context) Value(key interface{}) interface{} {
	panic("implement me")
}
