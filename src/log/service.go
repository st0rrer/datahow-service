package log

import "sync"

type Service struct {
	ipAddresses map[string]struct{}
	rwMutex     sync.RWMutex
}

func NewService() *Service {
	return &Service{
		ipAddresses: make(map[string]struct{}),
	}
}

func (s *Service) ProcessMessage(msg *Message) error {

	defer s.rwMutex.Unlock()
	s.rwMutex.Lock()
	s.ipAddresses[msg.IP] = struct{}{}

	return nil
}

func (s *Service) CountUniqueIP() (int, error) {

	defer s.rwMutex.RUnlock()
	s.rwMutex.RLock()

	return len(s.ipAddresses), nil
}
