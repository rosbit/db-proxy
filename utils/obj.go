package utils

import "sync"

var allObjs struct {
	sync.RWMutex
	objs map[uint32]interface{}
	next uint32
}

func init() {
	allObjs.objs = make(map[uint32]interface{})
	allObjs.next = 1
}

func NewObjId(obj interface{}) uint32 {
	allObjs.Lock()
	defer allObjs.Unlock()

	id := allObjs.next
	allObjs.next++
	if allObjs.next == 0 {
		allObjs.next = 1
	}
	allObjs.objs[id] = obj
	return id
}

func GetObjById(id uint32) interface{} {
	allObjs.RLock()
	defer allObjs.RUnlock()

	return allObjs.objs[id]
}

func FreeObjId(id uint32) interface{} {
	allObjs.Lock()
	defer allObjs.Unlock()

	obj := allObjs.objs[id]
	delete(allObjs.objs, id)

	return obj
}
