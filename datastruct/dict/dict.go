package dict

type Consumer func(key string, val interface{}) bool
type Dict interface {
	Get(key string) (val interface{}, exist bool)
	Set(key string, val string)
	Len() int
	Put(key string, val interface{}) (result int)
	PutIfAbsent(key string, val interface{}) (result int)
	PutIfExists(key string, val interface{}) (result int)
	Remove(key string) (result int)
	Foreach(consumer Consumer)
	Keys() []string
	RandomKeys(limit int) []string
	RandomDistinctKeys(limit int) []string
	Clear()
}
