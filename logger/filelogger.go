package logger

import (
	"bufio"
	"fmt"
	"os"
)

type FileLogger struct {
	log         *os.File
	lastEventId IdType
	events      chan<- Event
	errors      <-chan error
}

func NewFileLogger(path string) (Logger, error) {
	log, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("can't open log file: %w", err)
	}
	return &FileLogger{log: log}, nil
}

func (l *FileLogger) WritePut(key, val string) {
	l.events <- Event{
		Type: EventPut,
		Key:  key,
		Val:  val,
	}
}

func (l *FileLogger) WriteDelete(key string) {
	l.events <- Event{
		Type: EventDelete,
		Key:  key,
	}
}

func (l *FileLogger) Err() <-chan error {
	return l.errors
}

// TODO: create a config
const eventChanCapacity = 16
const errChanCapacity = 1
const lineFormat = "%d\t%d\t%s\t%s"

func (l *FileLogger) Run() {
	events := make(chan Event, eventChanCapacity)
	l.events = events

	errors := make(chan error, errChanCapacity)
	l.errors = errors

	go func() {
		for e := range events {

			l.lastEventId++

			// quick fix: scanning won't work
			// if we simply provide an empty string
			// TODO: come up with better solution
			var val string
			if len(e.Val) == 0 {
				val = "placeholder"
			} else {
				val = e.Val
			}

			_, err := fmt.Fprintf(
				l.log,
				lineFormat+"\n",
				l.lastEventId, e.Type, e.Key, val,
				// l.lastEventId, e.Type, e.Key, e.Val,
			)

			if err != nil {
				errors <- err
				return
			}
		}
	}()
}

func (l *FileLogger) ReadEvents() (<-chan Event, <-chan error) {
	scanner := bufio.NewScanner(l.log)
	outEvent := make(chan Event)
	outErr := make(chan error, errChanCapacity)

	go func() {
		defer close(outEvent)
		defer close(outErr)

		var e Event // read values will go here

		for scanner.Scan() {
			line := scanner.Text()

			if _, err := fmt.Sscanf(
				line,
				lineFormat,
				&e.Id, &e.Type, &e.Key, &e.Val,
			); err != nil {
				outErr <- fmt.Errorf("error occured during reading log file: %w", err)
				return
			}

			// integrity check
			if l.lastEventId >= e.Id {
				outErr <- fmt.Errorf("log ids out of order")
				return
			}

			l.lastEventId = e.Id

			outEvent <- e
		}

		if err := scanner.Err(); err != nil {
			outErr <- fmt.Errorf("error occured when trying to scan log: %w", err)
			return
		}
	}()

	return outEvent, outErr
}
