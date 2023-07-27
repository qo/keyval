package store

import (
	"errors"
	"fmt"
	"sync"

	"github.com/qo/keyval/logger"
)

// Signal errors
var (
	ErrNoSuchKey = errors.New("no such key")
	ErrEmptyKey  = errors.New("empty keys are not allowed")
	ErrEmptyVal  = errors.New("empty val is now allowed")
)

// Key-value store
type Store struct {
	sync.RWMutex
	m map[string]string
}

func CreateStore() (*Store, error) {
	return &Store{m: make(map[string]string)}, nil
}

func (s *Store) Put(key string, val string) error {
	if len(key) == 0 {
		return ErrEmptyKey
	}
	if len(val) == 0 {
		return ErrEmptyVal
	}
	s.Lock()
	s.m[key] = val
	s.Unlock()
	return nil
}

func (s *Store) Get(key string) (string, error) {
	if len(key) == 0 {
		return "", ErrEmptyKey
	}
	s.RLock()
	val, ok := s.m[key]
	s.RUnlock()
	if !ok {
		return "", ErrNoSuchKey
	}
	if len(val) == 0 {
		return "", ErrEmptyVal
	}
	return val, nil
}

func (s *Store) Delete(key string) error {
	if len(key) == 0 {
		return ErrEmptyKey
	}
	s.Lock()
	delete(s.m, key)
	s.Unlock()
	return nil
}

func (s *Store) InitLogger() (logger.Logger, error) {

	// TODO: create config file
	l, err := logger.NewSQLiteLogger(logger.SQLiteConfig{
		Path: "log.db",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init logger: %w", err)
	}

	// Read previous events from logger to
	// restore store state
	events, errors := l.ReadEvents()

	e, ok := logger.Event{}, true
	// while events and errors channels are not closed
	// and no operations caused errors
	for ok && err == nil {
		select {
		case err, ok = <-errors:
		case e, ok = <-events:
			switch e.Type {
			case logger.EventPut:
				err = s.Put(e.Key, e.Val)
			case logger.EventDelete:
				err = s.Delete(e.Key)
			}
		}
	}

	l.Run()

	return l, nil
}
