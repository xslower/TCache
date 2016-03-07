package main

import (
	"sync"
)

type Unit struct {
	Val   int64
	Score float32
}

func NewSet() *Set {
	var s = &Set{}
	s.Ints = make([]Unit, INTS_STEP)
	return s
}

type Set struct {
	Ints  []Unit
	RwMtx sync.RWMutex
	Head  uint32
	Tail  uint32
}

func (this *Set) Enlarge(count uint32) bool {
	this.RwMtx.Lock()
	defer this.RwMtx.Unlock()
	var ints = this.Ints
	if len(ints) >= int(_max_list_length) {
		return false
	}
	this.Ints = make([]Unit, len(ints)+int(count))
	copy(this.Ints, ints[this.Head:this.Tail])
	this.Head = 0
	this.Tail -= this.Head
	return true
}

func (this *Set) Shrink() bool {
	var step_int, step_ui32 = int(INTS_STEP), INTS_STEP
	this.RwMtx.Lock()
	defer this.RwMtx.Unlock()
	var ints = this.Ints
	if cap(ints) <= step_int {
		return false
	}
	this.Ints = make([]Unit, cap(ints)-step_int)
	if step_ui32 > this.Tail {
		this.Head = 0
		this.Tail = 0
		return true
	}
	if this.Head > step_ui32 {
		step_ui32 = this.Head
	}
	copy(this.Ints, ints[step_ui32:this.Tail])
	this.Head = 0
	this.Tail -= step_ui32
	return true
}

func (this *Set) FetchAll() []Unit {
	this.RwMtx.RLock()
	var ret = this.Ints[this.Head:this.Tail]
	this.RwMtx.RUnlock()
	return ret
}

func (this *Set) FetchRange(start, end int) []Unit {
	this.RwMtx.RLock()
	var ints = this.Ints[this.Head:this.Tail]
	var ln = len(ints)
	if end == 0 || end > ln {
		end = ln
	}
	var ret = []Unit{}
	if start < end {
		ret = ints[start:end]
	}
	this.RwMtx.RUnlock()
	return ret
}

func (this *Set) FetchOne(idx int) Unit {
	var ret = this.FetchRange(idx, idx+1)
	return ret[0]
}

func (this *Set) PopLeftOne() (Unit, error) {
	if this.IsEmpty() {
		return Unit{}, err_no_data
	}
	this.RwMtx.Lock()
	var ret = this.Ints[this.Head]
	this.Head++
	if this.Head >= uint32(CLEAN_UP) {
		this.moveToHead(this.Head)
	}
	this.RwMtx.Unlock()
	return ret, nil
}

func (this *Set) PopLeftMulti(count uint32) []Unit {
	if this.IsEmpty() {
		return []Unit{}
	}
	this.RwMtx.Lock()
	var tail = this.Head + count
	if tail > this.Tail {
		tail = this.Tail
	}
	var ret = this.Ints[this.Head:tail]
	this.Head = tail
	if this.Head >= uint32(CLEAN_UP) {
		this.moveToHead(this.Head)
	}
	this.RwMtx.Unlock()
	return ret
}

func (this *Set) PopRightOne() (Unit, error) {
	if this.IsEmpty() {
		return Unit{}, err_no_data
	}
	this.RwMtx.Lock()
	this.Tail--
	var ret = this.Ints[this.Tail]
	this.RwMtx.Unlock()
	return ret, nil
}

func (this *Set) PopRightMulti(count uint32) []Unit {
	if this.IsEmpty() {
		return []Unit{}
	}
	this.RwMtx.Lock()
	var head = this.Tail - count
	if head < this.Head {
		head = this.Head
	}
	var ret = this.Ints[head:this.Tail]
	this.Tail = head
	this.RwMtx.Unlock()
	return ret
}

func (this *Set) Push(val Unit) {
	this.RwMtx.Lock()
	this.Ints[this.Tail] = val
	this.Tail++
	this.RwMtx.Unlock()
}

func (this *Set) PushMulti(vals []Unit) {
	this.RwMtx.Lock()
	var ln = uint32(len(vals))
	var count = this.WantSomeSpace(ln)
	if count < ln { //如果容量不够则抛弃前面的值
		vals = vals[ln-count:]
	}
	copy(this.Ints[this.Tail:], vals)
	this.Tail += count
	this.RwMtx.Unlock()
}

func (this *Set) InsertAt(val Unit, pos uint32) {
	this.RwMtx.Lock()
	this.Tail++
	var ints = this.Ints[this.Head:this.Tail]
	copy(ints[pos+1:], ints[pos:])
	ints[pos] = val
	this.RwMtx.Unlock()
}

func (this *Set) WantSomeSpace(count uint32) uint32 {
	var ln = uint32(len(this.Ints))
	var tail = this.Tail + count
	if tail > ln {
		if count < CLEAN_UP {
			count = CLEAN_UP
		} else if count > ln {
			count = ln
		}
		this.moveToHead(count)
	}
	return count

}

/**
 * 把slice中start后面的数据整体移动头部，主要用来删除数据的
 */
func (this *Set) moveToHead(start uint32) {
	if start >= this.Tail { //mean clean all data
		this.Head = 0
		this.Tail = 0
		return
	}
	if this.Head > start {
		start = this.Head
	}
	copy(this.Ints, this.Ints[start:this.Tail])
	this.Head = 0
	this.Tail -= start
}

func (this *Set) IsFull() bool {
	this.RwMtx.RLock()
	var ret = int(this.Tail) >= len(this.Ints)
	this.RwMtx.RUnlock()
	return ret
}

func (this *Set) IsEmpty() bool {
	this.RwMtx.RLock()
	var ret = this.Head == this.Tail
	this.RwMtx.RUnlock()
	return ret
}

func (this *Set) Cap() int {
	return cap(this.Ints)
}

//bottom
//
