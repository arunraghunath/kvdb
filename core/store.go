package core

import (
	"time"
)

var store map[string]*KVObj

type KVObj struct {
	Value     string
	ExpiresAt int64
}

func init() {
	store = make(map[string]*KVObj)
}

func NewKVObj(val string, durationMs int64) *KVObj {
	var lexpiresAt int64 = -1
	if durationMs > 0 {
		lexpiresAt = time.Now().UnixMilli() + durationMs
	}

	return &KVObj{
		Value:     val,
		ExpiresAt: lexpiresAt,
	}
}

func Put(key string, kvobj *KVObj) {
	store[key] = kvobj
}

func Get(key string) *KVObj {
	return store[key]
}
