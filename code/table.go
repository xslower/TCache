package main

import (
	"sort"
)

const (
	LCL_MEM_STEP uint32 = 100
	INTS_STEP    uint32 = 200
	CLEAN_UP     uint32 = 100
	SEARCH_DEPTH uint32 = 200
	MIN_LEN      uint32 = 100
	HASH_RATE    uint32 = 4
	HOLE_SLOT    uint32 = 200
)

var (
	UI32_FALSE uint32 = 0
)

func initDb(max_mem, db_num uint32) {
	UI32_FALSE--
	var each_part = int(max_mem / db_num)
	_all_tables = make([]*AllTable, db_num)
	for i := uint32(0); i < db_num; i++ {
		var m = NewMemory(each_part)
		_all_tables[i] = NewAllTable(m)
	}
}

func NewAllTable(m *Memory) *AllTable {
	var at = &AllTable{}
	at.it = NewIntTable(m)
	at.lt = NewListTable(m, false)
	at.st = NewListTable(m, true)
	at.zt = NewSetTable(m)
	return at
}

type AllTable struct {
	it *IntTable
	lt *ListTable
	st *ListTable
	zt *SetTable
}

func (this *AllTable) GetIDel(t string) ITbDel {
	switch t {
	case `i`:
		return this.it
	case `l`:
		return this.lt
	case `s`:
		return this.st
	case `z`:
		return this.zt
	default:
		return nil
	}
}

func (this *AllTable) GetLT(t string) *ListTable {
	switch t {
	case `l`:
		return this.lt
	case `s`:
		return this.st
	default:
		return nil
	}
}

type ITbSet interface {
	Set(uint32, int64)
}

type ITbDel interface {
	Del(uint32)
}

func NewHole() *Hole {
	var h = &Hole{}
	h.holes = make([]uint32, HOLE_SLOT)
	return h
}

type Hole struct {
	holes []uint32
	tail  uint32
}

func (this *Hole) Len() uint32 {
	return this.tail
}
func (this *Hole) IsFull() bool {
	return int(this.tail) >= len(this.holes)
}

/**
 * 找比kp大的数中最小的pos，插入到pos并返回其值，没找到则插在最后并返回kp
 */
func (this *Hole) Replace(kp uint32) uint32 {
	var pos = this.Search(kp)
	if pos == this.tail {
		this.holes[this.tail] = kp
		this.tail++
		return kp
	} else {
		var ret = this.holes[pos]
		this.holes[pos] = kp
		return ret
	}
}

func (this *Hole) Set(kp uint32) {
	var pos = this.Search(kp)
	this.insertAt(kp, pos)
}

func (this *Hole) insertAt(kp uint32, pos uint32) {
	if pos < this.tail {
		copy(this.holes[pos+1:], this.holes[pos:this.tail])
	}
	this.holes[pos] = kp
	this.tail++
}

func (this *Hole) ClearAll() []uint32 {
	var ret = this.holes[:this.tail]
	this.tail = 0
	return ret
}

func (this *Hole) Search(kp uint32) uint32 {
	var pos = sort.Search(int(this.tail), func(i int) bool { return this.holes[i] >= kp })
	return uint32(pos)
}

func (this *Hole) GetAll() []uint32 {
	return this.holes[:this.tail]
}

func (this *Hole) Pop(max uint32) uint32 {
	if this.tail == 0 {
		return max
	}
	var ret = this.holes[this.tail]
	this.tail--
	return ret
}

func (this *Hole) Move(dis uint32) {
	var pos = this.Search(dis)
	if pos < this.tail {
		copy(this.holes, this.holes[pos:this.tail])
		this.tail -= pos
	}
	for i := uint32(0); i < this.tail; i++ {
		this.holes[i] -= dis
	}
}

/**
 * bottom
 */
