package pomodoro

import "sync"

type InMemStore struct {
	sync.RWMutex
	Pomodoros []Config
}

func (st *InMemStore) Create (c Config) (int, error){
	st.Lock()
	defer st.Unlock()
	c.ID = len(st.Pomodoros)+1
	st.Pomodoros = append(st.Pomodoros, c)
	return c.ID, nil
}

