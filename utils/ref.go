package utils

import "sync"

type FnCreate func(params interface{})(obj interface{}, err error)
type FnRelease func(obj interface{})

type ref struct {
	objId uint32
	release FnRelease
	refCount int
}

var allRefs struct {
	sync.RWMutex
	refs map[interface{}]*ref
	objId2Ref map[uint32]interface{}
}

func init() {
	allRefs.refs = make(map[interface{}]*ref)
	allRefs.objId2Ref = make(map[uint32]interface{})
}

func IncRef(params interface{}, create FnCreate, release FnRelease) (objId uint32, err error) {
	allRefs.Lock()
	defer allRefs.Unlock()

	r, ok := allRefs.refs[params]
	if ok {
		r.refCount += 1
		objId = r.objId
		return
	}

	obj, e := create(params)
	if e != nil {
		err = e
		return
	}

	objId = NewObjId(obj)
	allRefs.refs[params] = &ref{
		objId: objId,
		release: release,
		refCount: 1,
	}
	allRefs.objId2Ref[objId] = params
	return
}

func GetRef(params interface{}) (obj interface{}) {
	allRefs.RLock()
	defer allRefs.RUnlock()

	r, ok := allRefs.refs[params]
	if !ok {
		return
	}
	obj = GetObjById(r.objId)
	return
}

func GetIdByRef(params interface{}) (objId uint32, ok bool) {
	allRefs.RLock()
	defer allRefs.RUnlock()

	var r *ref
	r, ok = allRefs.refs[params]
	if !ok {
		return
	}
	objId = r.objId
	return
}

func GetRefByObjId(objId uint32) (obj interface{}) {
	allRefs.RLock()
	defer allRefs.RUnlock()

	if _, ok := allRefs.objId2Ref[objId]; !ok {
		return
	}
	obj = GetObjById(objId)
	return
}

func DecRef(params interface{}) {
	allRefs.Lock()
	defer allRefs.Unlock()

	r, ok := allRefs.refs[params]
	if !ok {
		return
	}

	r.refCount -= 1
	if r.refCount <= 0 {
		obj := FreeObjId(r.objId)
		delete(allRefs.refs, params)
		delete(allRefs.objId2Ref, r.objId)
		r.release(obj)
	}
}

func DecRefByObjId(objId uint32) {
	allRefs.Lock()
	defer allRefs.Unlock()

	params, ok := allRefs.objId2Ref[objId]
	if !ok {
		return
	}

	r, ok := allRefs.refs[params]
	if !ok {
		return
	}

	r.refCount -= 1
	if r.refCount <= 0 {
		obj := FreeObjId(r.objId)
		delete(allRefs.refs, params)
		delete(allRefs.objId2Ref, r.objId)
		r.release(obj)
	}
}
