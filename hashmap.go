package seq

import (
	"fmt"
)

// Hash maps are built on top of hash sets. KeyVal implements Setable, but the
// Hash and Equal methods only apply to the key and ignore the value.

// Container for a key/value pair, used by HashMap to hold its data
type KV struct {
	Key interface{}
	Val interface{}
}

func KeyVal(key, val interface{}) *KV {
	return &KV{key, val}
}

// Implementation of Hash for Setable. Only actually hashes the Key field
func (kv *KV) Hash(i uint32) uint32 {
	return hash(kv.Key, i)
}

// Implementation of Equal for Setable. Only actually compares the key field. If
// compared to another KV, only compares the other key as well.
func (kv *KV) Equal(v interface{}) bool {
	if kv2, ok := v.(*KV); ok {
		return equal(kv.Key, kv2.Key)
	}
	return equal(kv.Key, v)
}

// Implementation of String for Stringer
func (kv *KV) String() string {
	return fmt.Sprintf("%v -> %v", kv.Key, kv.Val) 
}

// HashMaps are actually built on top of Sets, just with some added convenience
// methods for interacting with them as actual key/val stores
type HashMap struct {
	set *Set
}

// Returns a new HashMap of the given KVs (or possibly just an empty HashMap)
func NewHashMap(kvs ...*KV) *HashMap {
	ints := make([]interface{}, len(kvs))
	for i := range kvs {
		ints[i] = kvs[i]
	}
	return &HashMap{
		set: NewSet(ints...),
	}
}

// Implementation of FirstRest for Seq interface. First return value will
// always be a *KV or nil. Completes in O(log(N)) time.
func (hm *HashMap) FirstRest() (interface{}, Seq, bool) {
	if hm == nil {
		return nil, nil, false
	}
	el, nset, ok := hm.set.FirstRest()
	return el, &HashMap{nset.(*Set)}, ok
}

// Returns a new HashMap with the given value set on the given key. Also returns
// whether or not this was the first time setting that key (false if it was
// already there and was overwritten). Has the same complexity as Set's SetVal
// method.
func (hm *HashMap) Set(key, val interface{}) (*HashMap, bool) {
	if hm == nil {
		hm = NewHashMap()
	}

	nset, ok := hm.set.SetVal(KeyVal(key, val))
	return &HashMap{nset}, ok
}

// Returns a new HashMap with the given key removed from it. Also returns
// whether or not the key was already there (true if so, false if not). Has the
// same time complexity as Set's DelVal method.
func (hm *HashMap) Del(key interface{}) (*HashMap, bool) {
	if hm == nil {
		hm = NewHashMap()
	}

	nset, ok := hm.set.DelVal(KeyVal(key, nil))
	return &HashMap{nset}, ok
}

// Returns a value for a given key from the HashMap, along with a boolean
// indicating whether or not the value was found. Has the same time complexity
// as Set's GetVal method.
func (hm *HashMap) Get(key interface{}) (interface{}, bool) {
	if hm == nil {
		return nil, false
	} else if kv, ok := hm.set.GetVal(KeyVal(key, nil)); ok {
		return kv.(*KV).Val, true
	} else {
		return nil, false
	}
}

// Same as FirstRest, but returns values already casted, which may be convenient
// in some cases.
func (hm *HashMap) FirstRestKV() (*KV, *HashMap, bool) {
	if el, nhm, ok := hm.FirstRest(); ok {
		return el.(*KV), nhm.(*HashMap), true
	} else {
		return nil, nil, false
	}
}

// Implementation of String for Stringer interface
func (hm *HashMap) String() string {
	return ToString(hm, "{", "}")
}

// Returns the number of KVs in the HashMap. Has the same complexity as Set's
// Size method.
func (hm *HashMap) Size() uint64 {
	return hm.set.Size()
}
