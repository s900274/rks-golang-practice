package common

type Context struct {
	Caller string
	LogId  int64
	Event  string
}

func NewContext(caller string, logid int64, event string) *Context {
	return &Context{
		Caller: caller,
		LogId:  logid,
		Event:  event,
	}
}
