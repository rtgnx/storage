package kv
type KVStore interface {
	Put(string, interface{}) error
	Get(string, interface{}) error
	Del(string) error
	Keys(string) ([]string, error)
	Exists(string) (bool, error)
}
