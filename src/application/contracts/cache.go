package contracts

type Cache[Key comparable, Value any] interface {
	Get(key Key) (value Value, ok bool)
	Set(key Key, value Value)
}
