package main

import (
	"sort"
)

func NewSortedSet() *SortedSet {
	var ss = &SortedSet{}
	ss.Set = *NewSet()
	return ss
}

type SortedSet struct {
	Set
}

func (this *SortedSet) Search(val float32, precise bool) int {
	this.RwMtx.RLock()
	var subject = this.Ints[this.Head:this.Tail]
	var f = func(i int) bool { return subject[i].Score == val }
	if !precise {
		f = func(i int) bool { return subject[i].Score >= val }
	}
	var pos = sort.Search(len(subject), f)
	if pos == len(subject) {
		pos = -1
	}
	this.RwMtx.RUnlock()
	return pos
}

func (this *SortedSet) FetchBetween(small, big float32) (ret []Unit) {
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

func (this *SortedSet) FetchBetweenLimit(small, big float32, limit int, odr Order) []Unit {
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
		var slc = make([]Unit, len(ret))
		copy(slc, ret)
		sort.Sort(sort.Reverse(MarkSlice(slc)))
		return slc
	}
	return ret
}

//不提供去重，麻烦而且意义不大
func (this *SortedSet) Push(s Unit) {
	var pos = this.Search(s.Score, false)
	if pos > -1 { //找到了
		var ts = this.FetchOne(pos)
		if ts.Score != s.Score || ts.Val != s.Val {
			this.Set.InsertAt(s, uint32(pos))
		} //Score与Val都相等的判定为重复元素，跳过
	} else { //没找到则说明val最大，所以放到最后
		this.Set.Push(s)
	}
}

func (this *SortedSet) Sort() {
	this.RwMtx.Lock()
	var ints = this.Ints[this.Head:this.Tail]
	sort.Sort(MarkSlice(ints))
	this.RwMtx.Unlock()
}

//bottom
