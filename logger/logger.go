package logger

type Logger interface {
	WritePut(key, val string)
	WriteDelete(key string)
	Err() <-chan error
	ReadEvents() (<-chan Event, <-chan error)
	Run()
}

type Event struct {
	Id   IdType
	Type EventType
	Key  string
	Val  string
}

type IdType uint64
type EventType int

const (
	EventPut EventType = iota + 1
	EventDelete
)
