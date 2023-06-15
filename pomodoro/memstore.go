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

func (st *InMemStore) Delete (id int) error {
	st.Lock()
	defer st.Unlock()
	var confs []Config
	if id == 0 {
		return fmt.Errorf("%w: %d", ErrInvalidID, id)
	}
	for _, v := range st.Pomodoros {
		if v.ID != id-1 {
			confs = append(confs, v)
		}
	}
	st.Pomodoros = confs
	return nil
}

func (st *InMemStore) GetCompleted () []Config {
	st.Lock()
	defer st.Unlock()
	var completed []Config
	for _, v := range st.Pomodoros {
		if v.State == StateDone {
			completed = append(completed, v)
		}
	}
	return completed
}

func (r *InMemStore) Breaks() ([]Config, error) {
	r.RLock()
	defer r.RUnlock()
	data := []Config{}
	for k := len(r.Pomodoros) - 1; k >= 0; k-- {
		if r.Pomodoros[k].Category == CatPomodoro {
			continue
		}
		data = append(data, r.Pomodoros[k])
	}
	return data, nil
}