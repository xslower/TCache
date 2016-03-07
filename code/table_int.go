package main

import (
	"sync"
)

/**
 * 因为原来的循环查找版的性能太低，现在改用散列表保存。
 * 直接用key%len，如果冲突则往后面插.
 * 在冲突中间的key，访问就排到前面，这样在冲突的key之间就形成了lru
 */
func NewIntTable(m *Memory) *IntTable {
	var it = &IntTable{Mem: m}
	it.Keys = make([]uint32, MIN_LEN)
	it.Vals = make([]int64, MIN_LEN)
	it.left = MIN_LEN
	return it
}

type IntTable struct {
	Keys  []uint32
	Vals  []int64
	Mem   *Memory
	mtx   sync.Mutex
	left  uint32 //剩余的slot数
	Visit uint16 //访问的次数，用来调节内存的
}

func (this *IntTable) enlarge(count int) {
	var tk = this.Keys
	var tv = this.Vals
	var org = len(tk)
	this.Keys = make([]uint32, org+count)
	this.Vals = make([]int64, org+count)
	copy(this.Keys, tk[:org])
	copy(this.Vals, tv[:org])
	this.tidyUp(uint32(org))

}

//扩容之后需要把原来的值重新分配
func (this *IntTable) tidyUp(org uint32) {
	var key uint32
	var val int64
	for i := uint32(0); i < org; i++ {
		if this.Keys[i] > org { //只有大于org的才会hash到后面
			this.left++
			key, val = this.Keys[i], this.Vals[i]
			this.Keys[i] = 0
			this.put(key, val)
		}
	}
}

func (this *IntTable) askForMem() {
	var step = getNextStep(len(this.Keys))
	//以4为单位
	if this.Mem.AskFor(step * 4) {
		this.enlarge(step)
		this.left += uint32(step)
	} //申请不到则do nothing
}

func (this *IntTable) leftSpace() uint32 {
	return this.left
}

func (this *IntTable) hashPos(key uint32) uint32 {
	return key % uint32(len(this.Keys))
}

func (this *IntTable) search(key uint32) uint32 {
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

/**
 * 把命中的元素从冲突链表的后面移动第一
 */
func (this *IntTable) visit(pos uint32) {
	this.Visit++
	var key, val = this.Keys[pos], this.Vals[pos]
	var ln = uint32(len(this.Keys))
	var desPos = key % ln
	var p = desPos
	for p != pos || this.Keys[p]%ln != desPos { //hash
		p++
		if p == ln {
			p = 0
		}
	}
	if p == pos {
		return
	}
	if p < pos {
		copy(this.Keys[p+1:pos+1], this.Keys[p:pos])
		copy(this.Vals[p+1:pos+1], this.Vals[p:pos])
	} else { //冲突到Keys的底部又循环到头部去了，此时移动内存太麻烦，直接插入
		//TODO
	}
	this.Keys[p], this.Vals[p] = key, val

}

func (this *IntTable) Get(key uint32) int64 {
	this.mtx.Lock()
	var ret int64 = 0
	var pos = this.search(key)
	if this.Keys[pos] == key {
		ret = this.Vals[pos]
		this.visit(pos)
	}
	this.mtx.Unlock()
	return ret
}

func (this *IntTable) put(key uint32, val int64) {
	var ln = uint32(len(this.Keys))
	var pos = this.search(key)
	if this.Keys[pos] == key { //有这个key则直接赋值即可
		this.Vals[pos] = val
	} else if this.Keys[pos] == 0 { //空slot插入
		this.Keys[pos] = key
		this.Vals[pos] = val
		this.left--
	} else {
		/**
		 * 附近都满了，有两种情况：
		 * 1.目标位置被前面占了，则抢前面最后一个位置，
		 * 2.目标位置被本链条兄弟占了，则抢本身冲突链条的最后一个位置，
		 */
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
		this.Keys[p-1], this.Vals[p-1] = key, val
	}
}

func (this *IntTable) Set(key uint32, val int64) {
	this.mtx.Lock()
	var ln = uint32(len(this.Keys))
	if this.left < ln/HASH_RATE {
		this.askForMem()
	}
	this.put(key, val)
	this.mtx.Unlock()
}

func (this *IntTable) Del(key uint32) {
	this.mtx.Lock()
	var pos = this.search(key)
	if this.Keys[pos] == key {
		this.Keys[pos] = 0
	}
	this.mtx.Unlock()
}

//num可以为负，所以不需要decrease
func (this *IntTable) Increase(key uint32, num int64) int64 {
	this.mtx.Lock()
	var pos = this.search(key)
	var val int64
	if this.Keys[pos] == key {
		val = this.Vals[pos]
		val += num
		this.Vals[pos] = val
	}
	this.mtx.Unlock()
	return val
}
