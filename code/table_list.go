package main

import (
	"sync"
)

func _newSList() IList {
	return NewSortedList()
}
func _newList() IList {
	return NewList()
}

/**
 * 有序表不比无序比更占内存，所以用有序表代替无序表和有序表
 */
func NewListTable(m *Memory, s bool) *ListTable {
	var lt = &ListTable{Mem: m}
	lt.Keys = make([]uint32, MIN_LEN)
	lt.Vals = make([]IList, MIN_LEN)
	lt.lru = *NewLRU()
	if s {
		lt.newList = _newSList
	} else {
		lt.newList = _newList
	}
	return lt
}

type ListTable struct {
	Keys    []uint32
	Vals    []IList
	Mem     *Memory
	lru     LRU
	mtx     sync.Mutex
	newList func() IList
	left    uint32
	lclMem  uint32 //本地内存池
	Visit   uint16 //访问的次数，用来调节内存
}

func (this *ListTable) enlarge(count int) {
	var tk = this.Keys
	var tv = this.Vals
	var org = len(tk)
	this.Keys = make([]uint32, org+count)
	this.Vals = make([]IList, org+count)
	copy(this.Keys, tk[:org])
	copy(this.Vals, tv[:org])
}

func (this *ListTable) tidyUp(org uint32) {
	var key uint32
	var val IList
	var pos uint32
	for i := uint32(0); i < org; i++ {
		if this.Keys[i] > org { //只有大于org的才会hash到后面
			this.left++
			key, val = this.Keys[i], this.Vals[i]
			this.Keys[i] = 0
			this.Vals[i] = nil
			pos = this.search(key)
			this.put(key, val, pos)
		}
	}
}

/**
 * 为自身table长度申请内存
 */
func (this *ListTable) askForMem() {
	var step = getNextStep(len(this.Keys))
	//以4为单位, 因为还有lru占内存
	if this.Mem.AskFor(step * 4) {
		this.enlarge(step)
		this.left += uint32(step)
	} //申请不到do nothing

}

func (this *ListTable) freeMem(count uint32) uint32 {
	var pos uint32
	var total int
	var dels = this.lru.Del(count)
	for _, key := range dels {
		pos = this.search(key)
		if this.Keys[key] == key {
			total += this.Vals[pos].Cap()
			this.Keys[pos] = 0
			this.Vals[pos] = nil //这里必须设为nil，不然内存无法释放
		}
	}
	return uint32(total)
}

func (this *ListTable) askMemForSub() {
	var count = LCL_MEM_STEP
	//list以INTS_STEP单位增加, 这里以2为单位
	var size = count * INTS_STEP
	if this.Mem.AskFor(int(size * 2)) {
		this.lclMem += size
	} else { //申请不到内存则删除前面的key空出空间
		var total = this.freeMem(count)
		this.lclMem += total
	}
}

func (this *ListTable) leftSpace() uint32 {
	return this.left
}

func (this *ListTable) search(key uint32) uint32 {
	var ln = uint32(len(this.Keys))
	var pos = key % ln //hash
	var i uint32 = 0
	for this.Keys[pos] != key && this.Keys[pos] != 0 {
		pos++
		if pos == ln {
			pos = 0
		}
		i++
		if i >= SEARCH_DEPTH {
			break
		}
	}
	return pos
}

func (this *ListTable) Get(key uint32) IList {
	this.mtx.Lock()
	var ret IList = nil
	var pos = this.search(key)
	if this.Keys[pos] == key {
		ret = this.Vals[pos]
		this.lru.Visit(key)
		this.Visit++
	}
	this.mtx.Unlock()
	return ret
}

func (this *ListTable) put(key uint32, val IList, pos uint32) {
	if this.Keys[pos] == 0 { //空slot插入
		this.Keys[pos] = key
		this.Vals[pos] = val
		this.left--
	} else {
		/**
		 * 附近都满了，有两种情况：
		 * 1.目标位置被前面占了，则抢前面最后一个位置，
		 * 2.目标位置被本链条兄弟占了，则抢本身冲突链条的最后一个位置，
		 */
		var ln = uint32(len(this.Keys))
		var p = key % ln
		var hashDes = this.Keys[p] % ln
		var hash = hashDes
		for hash == hashDes {
			p++
			if p == ln {
				p = 0
			}
			hash = this.Keys[p] % ln
		}
		p--
		this.lclMem += uint32(this.Vals[p].Cap()) //还回内存
		this.Keys[p], this.Vals[p] = key, val
	}
}

/**
 * 1.先判断内存池剩余内存是否足够，不够则申请内存或删除数据释放内存，同时如果slot满了或者enlarge或者删除数据释放slot
 * 2.
 */
func (this *ListTable) Set(key uint32, val int64) {
	this.mtx.Lock()
	var ln = uint32(len(this.Keys))
	if this.left < ln/HASH_RATE {
		this.askForMem()
	}
	if this.lclMem < INTS_STEP {
		this.askMemForSub()
	}
	var obj IList = nil
	var pos = this.search(key)
	if this.Keys[pos] == key { //有这个key则判断是否已满，如enlarge成功则减内存池内存
		obj = this.Vals[pos]
		if obj.IsFull() {
			if obj.Enlarge(INTS_STEP) {
				this.lclMem -= INTS_STEP
			} else { //扩容失败则清理空间
				_ = obj.WantSomeSpace(1)
			}
		}
	} else { //没有此key则新建list
		obj = this.newList()
		this.put(key, obj, pos)
		this.lclMem -= INTS_STEP
		this.lru.Add(key)
	}
	obj.Push(val)
	this.mtx.Unlock()
}

func (this *ListTable) Del(key uint32) {
	this.mtx.Lock()
	var pos = this.search(key)
	if this.Keys[pos] == key { //只有准确找到了此key此删除
		this.Keys[pos] = 0
		this.Vals[pos] = nil
	}
	this.mtx.Unlock()
}
