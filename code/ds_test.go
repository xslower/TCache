package main

import (
	"runtime"
	"runtime/debug"
	"testing"
	"time"
)

var (
	_m  = NewMemory(128)
	_at = NewAllTable(_m)
)

func init() {
	echo(`DS_TEST:`)
}

func TestSetGet(t *testing.T) {
	echo(`_SetGet:`)
	return
	var tb = _at.it
	for i := uint32(1); i < 99000; i++ {

		tb.Set(i, int64(i))
	}
	var val = tb.Get(12345)
	echo(val)
	runtime.GC()
	debug.FreeOSMemory()
	echo(len(tb.Keys), len(tb.Vals), tb.left)
	//	time.Sleep(time.Minute)
}

func TestList(t *testing.T) {
	echo(`_List:`)
	return
	var tb = _at.st

	tb.Set(1, 1)
	var start = time.Now().UnixNano()
	for i := uint32(1); i < 10001; i++ {
		for j := int64(1); j < 100; j++ {
			//			if j == 500 {
			//				_ = "breakpoint"
			//			}
			tb.Set(i, j)
		}
	}
	runtime.GC()
	debug.FreeOSMemory()
	var end = time.Now().UnixNano()
	_ = "breakpoint"
	var list = tb.Get(12)
	echo(list.FetchAll())
	list = tb.Get(7)
	echo(list.FetchAll())
	list = tb.Get(123)
	//	echo(list.FetchAll())
	var tail = tb.lru.holes.tail
	echo(end-start, tb.lru.holes.holes[:tail])
	time.Sleep(time.Minute)
}

func TestSortedSet(t *testing.T) {
	echo(`SortedSet:`)
	var tb = _at.zt
	var start = time.Now().UnixNano()
	var ut = Unit{1, 0.1}
	for i := uint32(1); i < 10001; i++ {
		for j := int64(1); j < 100; j++ {
			ut.Val = j
			tb.Set(i, ut)
		}
	}
	runtime.GC()
	debug.FreeOSMemory()
	var end = time.Now().UnixNano()
	var m1 = Unit{Val: 1234, Score: 10.1}
	tb.Set(12345, m1)
	_ = "breakpoint"
	var ss = tb.Get(12)
	var mk, err = ss.PopLeftOne()
	throw(err)
	if mk.Score != m1.Score || mk.Val != m1.Val {
		t.Error(`sorted set pop failed`)
	}
	echo(ss.FetchAll())
	var tail = tb.lru.holes.tail
	echo(end-start, tb.lru.holes.holes[:tail])
	time.Sleep(time.Minute)
}

func placeHolder() {
	_ = time.Minute
}

/**
bottom
*/
