package core

import (
	"fmt"
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

func Del(key string) int {
	lobj := Get(key)
	if lobj == nil {
		return 0
	}
	delete(store, key)
	return 1
}

func Expire() int {
	lcount := 0
	lcurr := time.Now().UnixMilli()
	fmt.Println(lcurr)
	for k, v := range store {
		fmt.Println("Key is ", k)
		fmt.Println("Value is ", v.Value, " ", v.ExpiresAt)
		if v.ExpiresAt != -1 && v.ExpiresAt < lcurr {
			fmt.Println("Deleting the key")
			lcount = lcount + Del(k)
		}
	}
	return lcount
}
