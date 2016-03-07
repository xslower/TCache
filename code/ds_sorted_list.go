package main

import (
	"sort"
)

type Order uint8

const (
	ORDER_ASC  Order = 1
	ORDER_DESC Order = 2
)

type IList interface {
	Push(int64)
	PopRightMulti(uint32) []int64
	PopLeftMulti(uint32) []int64
	FetchAll() []int64
	FetchBetweenLimit(int64, int64, int, Order) []int64
	Cap() int
	IsFull() bool
	Enlarge(uint32) bool
	WantSomeSpace(uint32) uint32
}

func NewSortedList() *SortedList {
	var sl = &SortedList{}
	sl.List = *NewList()
	return sl
}

type SortedList struct {
	List
}

func (this *SortedList) Search(val int64, precise bool) int {
	this.RwMtx.RLock()
	var subject = this.Ints[this.Head:this.Tail]
	var f = func(i int) bool { return subject[i] == val }
	if !precise {
		f = func(i int) bool { return subject[i] >= val }
	}
	var pos = sort.Search(len(subject), f)
	if pos == len(subject) {
		pos = -1
	}
	this.RwMtx.RUnlock()
	return pos
}

func (this *SortedList) FetchBetween(small, big int64) (ret []int64) {
	if big <= small {
		return
	}
	var start = this.Search(small, false)
	var end = this.Search(big, false)
	if start == -1 {
		if end == -1 {
			//do nothing
		} else {
			ret = this.FetchRange(0, end)
		}
	} else {
		if end == -1 {
			ret = this.FetchRange(start, 0)
		} else {
			ret = this.FetchRange(start, end)
		}
	}
	return
}

func (this *SortedList) FetchBetweenLimit(small, big int64, limit int, odr Order) []int64 {
	var ret = this.FetchBetween(small, big)
	if limit < len(ret) {
		if odr == ORDER_ASC {
			ret = ret[:limit]
		} else {
			ret = ret[len(ret)-limit:]
		}
	}
	if odr == ORDER_DESC {
		//如果不新分配内存，排序会把原数据位置修改
		var slc = make([]int64, len(ret))
		copy(slc, ret)
		sort.Sort(sort.Reverse(Int64Slice(slc)))
		return slc
	}
	return ret
}

func (this *SortedList) Push(val int64) {
	//如果存在一样的值，则不插入
	var pos = this.Search(val, false)
	if pos > -1 { //找到了
		if this.FetchOne(pos) != val {
			this.List.InsertAt(val, uint32(pos))
		} //如果找到的值与val相同则跳过
	} else { //没找到则说明val最大，所以放到最后
		this.List.Push(val)
	}
}

func (this *SortedList) sortMultiPush(vals []int64) {
	for _, v := range vals {
		this.Push(v)
	}
}

func (this *SortedList) Sort() {
	this.RwMtx.Lock()
	var ints = this.Ints[this.Head:this.Tail]
	sort.Sort(Int64Slice(ints))
	this.RwMtx.Unlock()
}

//bottom
