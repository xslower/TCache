package main

import (
	`errors`
	"sync"
)

func NewUnit() *Unit {
	var this = &Unit{}
	return this
}

type Unit struct {
	ValArr [UNIT_COUNT]int64
	Next   *Unit
	Prev   *Unit
	RwLock sync.RWMutex
	Head   int8
	Tail   int8
}

func (this *Unit) GetIndexLeft(idx int) int {
	this.RwLock.RLock()
	var val = this.Tail - this.Head
	this.RwLock.RUnlock()
	idx -= int(val)
	return idx
}

func (this *Unit) IsFull() bool {
	this.RwLock.RLock()
	var tail = this.Tail
	var head = this.Head
	this.RwLock.RUnlock()
	if tail == UNIT_COUNT && head == 0 {
		return true
	}
	return false
}

func (this *Unit) IsTailFull() bool {
	this.RwLock.RLock()
	var val = this.Tail
	this.RwLock.RUnlock()
	if val == UNIT_COUNT {
		return true
	}
	return false
}

func (this *Unit) IsHeadFull() bool {
	this.RwLock.RLock()
	var val = this.Head
	this.RwLock.RUnlock()
	if val == 0 {
		return true
	}
	return false
}

func (this *Unit) IsEmpty() bool {
	this.RwLock.RLock()
	var l = this.Head
	var r = this.Tail
	this.RwLock.RUnlock()
	if l == r {
		return true
	}
	return false
}

func (this *Unit) MoveToHead() {
	if this.IsHeadFull() {
		return
	}
	this.RwLock.Lock()
	var i int8 = 0
	for j := this.Head; j < this.Tail; j++ {
		this.ValArr[i] = this.ValArr[j]
		i++
	}
	this.Head = 0
	this.Tail = i
	this.RwLock.Unlock()
}

func (this *Unit) MoveToTail() {
	if this.IsTailFull() {
		return
	}
	this.RwLock.Lock()
	var i int8 = UNIT_COUNT - 1
	for j := this.Tail - 1; j >= this.Head; j-- {
		this.ValArr[i] = this.ValArr[j]
		i--
	}
	this.Head = i + 1
	this.Tail = UNIT_COUNT
	this.RwLock.Unlock()
}

func (this *Unit) FetchOne(idx int) (int64, error) {
	if this.GetIndexLeft(idx) > 0 {
		return 0, err_out_range
	}
	this.RwLock.RLock()
	var ret = this.ValArr[int(this.Head)+idx]
	this.RwLock.RUnlock()
	return ret, nil
}

func (this *Unit) FetchAll() []int64 {
	this.RwLock.RLock()
	var ret = this.ValArr[this.Head:this.Tail]
	this.RwLock.RUnlock()
	return ret
}

func (this *Unit) FetchRange(start, end int) []int64 {
	if this.GetIndexLeft(start) > 0 || this.GetIndexLeft(end) > 0 || start <= end {
		return []int64{}
	}
	this.RwLock.RLock()
	var ret = this.ValArr[int(this.Head)+start : int(this.Head)+end]
	this.RwLock.RUnlock()
	return ret
}

func (this *Unit) PopOne() (int64, error) {
	if this.IsEmpty() {
		return 0, err_no_data
	}
	this.RwLock.Lock()
	var ret = this.ValArr[this.Head]
	this.Head++
	this.RwLock.Unlock()
	return ret, nil
}

func (this *Unit) PopMulti(num int8) (ret []int64, err error) {
	this.RwLock.Lock()
	if num > this.Tail-this.Head {
		ret = this.ValArr[this.Head:this.Tail]
		err = err_no_enough
		this.Head = this.Tail
	} else {
		ret = this.ValArr[this.Head : this.Head+num]
		this.Head += num
	}
	this.RwLock.Unlock()
	return
}

func (this *Unit) PopAll() []int64 {
	this.RwLock.Lock()
	var ret = this.ValArr[this.Head:this.Tail]
	this.Head = this.Tail
	this.RwLock.Unlock()
	return ret
}

func (this *Unit) PushAtTail(val int64) bool {
	if this.IsTailFull() {
		if this.IsHeadFull() {
			return false
		} else {
			this.MoveToHead()
		}
	}
	this.RwLock.Lock()
	this.ValArr[this.Tail] = val
	this.Tail++
	this.RwLock.Unlock()
	return true
}

func (this *Unit) PushAtHead(val int64) bool {
	if this.IsHeadFull() {
		if this.IsTailFull() {
			return false
		} else {
			this.MoveToTail()
		}
	}
	this.RwLock.Lock()
	this.Head--
	this.ValArr[this.Head] = val
	this.RwLock.Unlock()
	return true
}

//TODO
func (this *Unit) InsertAt(val int64, pos int8) bool {
	if this.GetIndexLeft(int(pos)) > 0 {
		return false
	}
	// var d int8 = 1
	// defer this.RwLock.Unlock()
	// this.RwLock.Lock()
	// if this.Tail == UNIT_COUNT {
	// 	if this.Head == 0 {
	// 		return false
	// 	}
	// 	d = 0
	// }
	this.RwLock.Lock()
	this.ValArr[this.Head+pos] = val
	this.RwLock.Unlock()
	return true
}

func NewFreeUnits(count int) *FreeUnits {
	var this = &FreeUnits{}
	this.Count = count
	this.UnitPtr = NewUnit()
	var ptr = this.UnitPtr
	for i := 1; i < count; i++ {
		ptr.Next = NewUnit()
		ptr = ptr.Next
	}
	return this
}

type FreeUnits struct {
	UnitPtr *Unit
	Count   int
	Lock    sync.Mutex
}

func (this *FreeUnits) GetCount() int {
	return this.Count
}

func (this *FreeUnits) GetUnit() *Unit {
	if this.Count == 0 {
		return nil
	}
	this.Lock.Lock()
	var ret = this.UnitPtr
	this.UnitPtr = ret.Next
	this.Count--
	this.Lock.Unlock()
	return ret
}

func (this *FreeUnits) SetUnit(ut *Unit) {
	this.Lock.Lock()
	this.Count++
	ut.Next = this.UnitPtr
	this.UnitPtr = ut
	this.Lock.Unlock()
}
