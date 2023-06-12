package pomodoro

import (
	"fmt"
	"sync"
)

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

func (st *InMemStore) Update (c Config) error {
	st.Lock()
	defer st.Unlock()
	if c.ID == 0 {
		return fmt.Errorf("%w: %d", ErrInvalidID, c.ID)
	}
	st.Pomodoros[c.ID-1] = c
	return nil
}

func (st *InMemStore) GetById (id int) (Config, error) {
	st.Lock()
	defer st.Unlock()
	if id < 0 {
		return Config{}, fmt.Errorf("%w: %d", ErrInvalidID, id)
	}
	return st.Pomodoros[id-1], nil
}

