package session

import (
	"container/list"
	"sync"
)

const (
	MAX_STORED = 500
)


type lossytable struct{
	lookuplock sync.Mutex
	expiryQueue *list.List
	Lookup map[string]*session
}

type session struct{
	UID int64
	key string
	attribs map[string]string
	listpos *list.Element
}


func newSessionTable()*lossytable{
	return &lossytable{
		expiryQueue: list.New(),
		Lookup: make(map[string]*session),
	}
}




func (inst *lossytable)getSession(key string)*session{
	inst.lookuplock.Lock()
	ret := inst.Lookup[key]
	inst.lookuplock.Unlock()
	if ret != nil{
		inst.expiryQueue.MoveToFront(ret.listpos)
	}
	return ret
}


func (inst *lossytable)updateAttrs(Key, k, v string){
	s := inst.getSession(Key)
	s.attribs[k] = v
}


func (inst *lossytable)insert(useruid int64, Key string){
	inst.lookuplock.Lock()
	defer inst.lookuplock.Unlock()
	
	entry := &session{
		UID: useruid,
		key: Key,
		attribs: make(map[string]string),
	}
	inst.Lookup[Key] = entry
	entry.listpos = inst.expiryQueue.PushFront(entry)
	inst.boundStorage()
}


func (inst *lossytable)delete(key string) {
	s := inst.getSession(key)
	if s != nil {
		inst.expiryQueue.Remove(s.listpos)
		delete(inst.Lookup, key)
	}
}


//called internally to check if our table has grown too large, and if so, to clean it up.
//MUST ALREADY have lock on the routing table when this is called.
func (inst *lossytable)boundStorage(){
	if inst.expiryQueue.Len() >= MAX_STORED{
		node := inst.expiryQueue.Back()
		inst.expiryQueue.Remove(node)
		inst.Lookup[node.Value.(*session).key].listpos = nil
		delete(inst.Lookup, node.Value.(*session).key)
	}
}
