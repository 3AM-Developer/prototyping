package state

import (
	"sync"

	"github.com/3AM-Developer/server-runner/internal/instance"
)

var AppState = Init()

type Status int

const (
	Running Status = iota
	Stopped
)

type State struct {
	mu       sync.Mutex
	instance *instance.Instance
	status   Status
}

func Init() *State {
	return &State{
		mu:       sync.Mutex{},
		instance: nil,
		status:   Stopped,
	}
}

// errors here
var (
	ErrorInstanceAlreadyStarted error
	ErrorInstanceAlreadyStopped error
	ErrorInstnaceInvalid        error
)

func (s *State) GetInstance() (*instance.Instance, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.instance != nil {
		// Super explicit copy here because I don't want
		// the pointer to __ever__ get passed around
		return s.instance.Copy(), true
	}

	return nil, false
}

func (s *State) RegisterInstance(i *instance.Instance) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !i.VerifyInstance() {
		return false, ErrorInstnaceInvalid
	}

	if s.instance != nil {
		return false, nil
	}

	s.instance = i
	return true, nil
}

func (s *State) UnregisterInstance() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.instance == nil {
		return false
	}

	if s.status == Running {
		return false
	}

	s.instance = nil
	return true
}

func (s *State) StartInstance() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status == Running {
		return ErrorInstanceAlreadyStarted
	}

	err := s.instance.Start()
	if err != nil {
		return err
	}

	s.status = Running
	return nil
}

func (s *State) StopInstance() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status == Stopped {
		return ErrorInstanceAlreadyStopped
	}

	err := s.instance.Stop()
	if err != nil {
		return err
	}

	s.status = Stopped
	return nil
}
