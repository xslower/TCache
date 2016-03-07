package main

import (
	`sync`
	"time"
)

func NewList() *List {
	var ut = freeList.GetUnit()
	if ut == nil {
		return nil
	}
	var this = &List{}
	this.AddUnitAtHead(ut)
	this.AddUnitAtTail(ut)
	this.Visit()
	return this
}

type List struct {
	HeadPtr   *Unit
	TailPtr   *Unit
	RwMtx     sync.RWMutex
	LastVisit int64
	Frequency int16
}

func (this *List) FreeUnit(ut *Unit) {
	if ut.Prev != nil {
		ut.Prev.Next = ut.Next
	}
	if ut.Next != nil {
		ut.Next.Prev = ut.Prev
	}
	freeList.SetUnit(ut)
}

func (this *List) FreeHeadUnit() {
	this.FreeUnit(this.HeadPtr)
}

func (this *List) FreeTailUnit() {
	this.FreeUnit(this.TailPtr)
}

func (this *List) AddUnitAtHead(ut *Unit) {
	ut.Next = this.HeadPtr
	ut.Prev = nil
	this.HeadPtr = ut
}

func (this *List) AddUnitAtTail(ut *Unit) {
	ut.Prev = this.TailPtr
	ut.Next = nil
	this.TailPtr = ut
}

func (this *List) Visit() {
	var last = this.LastVisit
	var now = time.Now().Unix()
	//上次使用距离本次使用比较近才增加频率
	if now-last < 3600*6 {
		this.Frequency++
	}
	this.LastVisit = time.Now().Unix()
}

//idx = 0 ~ N-1
func (this *List) FetchOne(idx int) (int64, error) {
	var ut = this.HeadPtr
	return ut.FetchOne(idx)

}

func (this *List) FetchAll() []int64 {
	this.Visit()
	var ret []int64
	var ptr = this.HeadPtr
	ret = ptr.ValArr[ptr.Head:ptr.Tail]
	for ptr != nil {
		ptr = ptr.Next
		ret = append(ret, ptr.ValArr[ptr.Head:ptr.Tail]...)
	}
	return ret
}

func (this *List) FetchRange(start, end int) []int64 {
	return []int64{}
}

//bottom
