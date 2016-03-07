//namespace Queue\DataStructure
//
package main

import (
	"sync"
)

func NewList() *List {
	var this = &List{}
	this.Ints = make([]int64, INTS_STEP)
	return this
}

type List struct {
	Ints  []int64
	RwMtx sync.RWMutex
	Head  uint32
	Tail  uint32
}

func (this *List) Enlarge(count uint32) bool {
	this.RwMtx.Lock()
	defer this.RwMtx.Unlock()
	var ints = this.Ints
	if len(ints) >= int(_max_list_length) {
		return false
	}
	this.Ints = make([]int64, len(ints)+int(count))
	copy(this.Ints, ints[this.Head:this.Tail])
	this.Head = 0
	this.Tail -= this.Head
	return true
}

func (this *List) Shrink() bool {
	var step_int, step_ui32 = int(INTS_STEP), INTS_STEP
	this.RwMtx.Lock()
	defer this.RwMtx.Unlock()
	var ints = this.Ints
	if cap(ints) <= step_int {
		return false
	}
	this.Ints = make([]int64, cap(ints)-step_int)
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

func (this *List) FetchAll() []int64 {
	this.RwMtx.RLock()
	var ret = this.Ints[this.Head:this.Tail]
	this.RwMtx.RUnlock()
	return ret
}

func (this *List) FetchRange(start, end int) []int64 {
	this.RwMtx.RLock()
	var ints = this.Ints[this.Head:this.Tail]
	var ln = len(ints)
	if end == 0 || end > ln {
		end = ln
	}
	var ret = []int64{}
	if start < end {
		ret = ints[start:end]
	}
	this.RwMtx.RUnlock()
	return ret
}

func (this *List) FetchOne(idx int) int64 {
	var ret = this.FetchRange(idx, idx+1)
	return ret[0]
}

func (this *List) FetchBetweenLimit(small, big int64, limit int, odr Order) []int64 { //for IList
	return []int64{}
}

func (this *List) PopLeftOne() (int64, error) {
	if this.IsEmpty() {
		return 0, err_no_data
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

func (this *List) PopLeftMulti(count uint32) []int64 {
	if this.IsEmpty() {
		return []int64{}
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

func (this *List) PopRightOne() (int64, error) {
	if this.IsEmpty() {
		return 0, err_no_data
	}
	this.RwMtx.Lock()
	this.Tail--
	var ret = this.Ints[this.Tail]
	this.RwMtx.Unlock()
	return ret, nil
}

func (this *List) PopRightMulti(count uint32) []int64 {
	if this.IsEmpty() {
		return []int64{}
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

func (this *List) Push(val int64) {
	this.RwMtx.Lock()
	this.Ints[this.Tail] = val
	this.Tail++
	this.RwMtx.Unlock()
}

func (this *List) PushMulti(vals []int64) {
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

func (this *List) InsertAt(val int64, pos uint32) {
	this.RwMtx.Lock()
	this.Tail++
	var ints = this.Ints[this.Head:this.Tail]
	copy(ints[pos+1:], ints[pos:])
	ints[pos] = val
	this.RwMtx.Unlock()
}

func (this *List) WantSomeSpace(count uint32) uint32 {
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
func (this *List) moveToHead(start uint32) {
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

func (this *List) IsFull() bool {
	this.RwMtx.RLock()
	var ret = int(this.Tail) >= len(this.Ints)
	this.RwMtx.RUnlock()
	return ret
}

func (this *List) IsEmpty() bool {
	this.RwMtx.RLock()
	var ret = this.Head == this.Tail
	this.RwMtx.RUnlock()
	return ret
}

func (this *List) Cap() int {
	return cap(this.Ints)
}

//bottom
//
