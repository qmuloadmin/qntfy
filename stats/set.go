package stats

import "sync"

// Sometimes, I really wish go had sets.
type set struct {
	lock *sync.Mutex
	m    map[string]bool
}

func (s *set) init() {
	s.lock = &sync.Mutex{}
	s.m = make(map[string]bool)
}

func (s set) add(item string) {
	s.m[item] = true
}

func (s set) safeAdd(item string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.add(item)
}

func (s set) contains(item string) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, has := s.m[item]; has {
		return true
	}
	return false
}

func (s set) items() []string {
	items := make([]string, 0, len(s.m))
	for item := range s.m {
		items = append(items, item)
	}
	return items
}
