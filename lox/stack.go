package lox

import (
	"errors"
	"sync"
)

// ScopeEntry ...
type ScopeEntry struct {
	defined bool
	used    bool
	token   *Token
}

// NewScopeStack constructor
func NewScopeStack() *ScopeStack {
	return &ScopeStack{s: []map[string]*ScopeEntry{}}
}

// ScopeStack for scopes
type ScopeStack struct {
	s  []map[string]*ScopeEntry
	mu sync.Mutex
}

func (s *ScopeStack) Push(v map[string]*ScopeEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.s = append(s.s, v)
}

func (s *ScopeStack) Pop() (map[string]*ScopeEntry, error) {
	l := len(s.s)
	if l == 0 {
		return nil, errors.New("empty stack")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	v, ns := s.s[l-1], s.s[:l-1]
	s.s = ns

	return v, nil
}

func (s *ScopeStack) Size() int {
	return len(s.s)
}

func (s *ScopeStack) Peek() (map[string]*ScopeEntry, error) {
	l := len(s.s)
	if l == 0 {
		return nil, errors.New("empty stack")
	}
	return s.s[l-1], nil
}

func (s *ScopeStack) Get(i uint) (map[string]*ScopeEntry, error) {
	l := uint(len(s.s))
	if i > l-1 {
		return nil, errors.New("index out of range")
	}

	return s.s[i], nil
}
