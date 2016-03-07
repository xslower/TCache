package main

// "sync"
// `time`

const (
	LRU_CLEAN  uint32 = 50
	RECENT_NUM uint32 = 100
)

func NewLRU() *LRU {
	var lru = &LRU{holes: *NewHole()}
	lru.lru = make([]uint32, MIN_LEN)
	return lru
}

type LRU struct {
	lru   []uint32
	holes Hole
	tail  uint32
}

func (this *LRU) Enlarge() {
	var step = getNextStep(len(this.lru))
	var lr = this.lru
	this.lru = make([]uint32, cap(lr)+step)
	copy(this.lru, lr)
}

func (this *LRU) Reset() {
	this.lru = make([]uint32, MIN_LEN)
}

/**
 * 根据holes中记录的洞的位置，通过循环移动内存把洞补上
 */
func (this *LRU) trimHole() uint32 {
	var hls = this.holes.ClearAll()
	var ln = len(hls)
	var dis uint32 = 1
	var head, tail uint32
	for i := 0; i < ln; i++ {
		head = hls[i] + 1
		if i+1 == ln {
			tail = this.tail
		} else {
			tail = hls[i+1]
		}
		copy(this.lru[head-dis:tail], this.lru[head:tail])
		dis++
	}
	return uint32(ln)
}

func (this *LRU) leftSpace() uint32 {
	var ret = uint32(len(this.lru)) + this.holes.Len() - this.tail
	return ret
}

func (this *LRU) Add(key uint32) {
	if this.leftSpace() < 1 {
		this.Enlarge()
	}
	this.lru[this.tail] = key
	this.tail++
}

func (this *LRU) Del(tail uint32) (ret []uint32) {
	this.trimHole()
	if tail >= this.tail {
		ret = this.lru
		this.tail = 0
	} else {
		ret = this.lru[:tail]
		//这样lru长度会慢慢变短，让enlarge去增加吧
		this.lru = this.lru[tail:]
		this.tail -= tail
	}
	return
}

func (this *LRU) Visit(key uint32) {
	if this.holes.IsFull() {
		this.trimHole()
	}
	if this.leftSpace() < 1 {
		this.Enlarge()
	}
	var pos uint32
	for pos = this.tail - 1; pos >= 0; pos-- {
		if this.lru[pos] == key {
			break
		}
	}
	this.lru[pos] = 0
	var newPos = this.holes.Replace(pos)
	if newPos == pos { //说明没有
		newPos = this.tail
		this.tail++
	}
	this.lru[newPos] = key
}

//bottom
//
