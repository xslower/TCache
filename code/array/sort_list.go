package main

import (
	"time"
	// `sync`
)

func NewSortList() *SortList {
	var this = &SortList{}
	this.HeadPtr = freeList.GetUnit()
	this.HeadPtr.Next = freeList.GetUnit()
	this.TailPtr = this.HeadPtr.Next
	this.TailPtr.Prev = this.HeadPtr
	this.Visit()
	return this
}

type SortList struct {
	HeadPtr   *Unit
	TailPtr   *Unit
	LastVisit int64
	Frequency int16
	// Head    *int64
	// Tail    *int64
}

func (this *SortList) RemoveHeadUnit() {

}

func (this *SortList) RemoveTailUnit() {

}

func (this *SortList) Visit() {
	var last = this.LastVisit
	var now = time.Now().Unix()
	//上次使用距离本次使用比较近才增加频率
	if now-last < 3600*6 {
		this.Frequency++
	}
	this.LastVisit = time.Now().Unix()
}

func (this *SortList) Add(val int64) bool {

	return true
}

//bottom
