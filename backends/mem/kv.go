package memory

import (
	"encoding/json"
	"fmt"
	"path"
	"sync"

	"github.com/rtgnx/storage/kv"
)

type Memory struct {
	kv     map[string]string
	kvLock sync.Mutex
}

func New() (kv.KVStore, error) {
	return &Memory{kv: make(map[string]string)}, nil

}

func (m *Memory) Put(prefix string, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	m.kvLock.Lock()
	defer m.kvLock.Unlock()

	m.kv[prefix] = string(b)
	return nil
}

func (m *Memory) Get(prefix string, v interface{}) error {
	m.kvLock.Lock()
	defer m.kvLock.Unlock()

	if b, ok := m.kv[prefix]; ok {
		return json.Unmarshal([]byte(b), v)
	}

	return fmt.Errorf("key not found: %s", prefix)
}

func (m *Memory) Del(prefix string) error {
	m.kvLock.Lock()
	defer m.kvLock.Unlock()

	delete(m.kv, prefix)
	return nil
}

func (m *Memory) Exists(prefix string) (bool, error) {
	m.kvLock.Lock()
	defer m.kvLock.Unlock()

	_, ok := m.kv[prefix]
	return ok, nil
}

func (m *Memory) Keys(prefix string) ([]string, error) {
	m.kvLock.Lock()
	defer m.kvLock.Unlock()

	out := []string{}

	for k := range m.kv {
		if m, _ := path.Match(prefix, k); m {
			out = append(out, k)
		}
	}

	return out, nil
}
