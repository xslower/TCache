package main

import (
	"sync"
)

/**
 * 因为以大量使用uint32，所以以4B为单位
 * 则在64位系统上的指针/接口/int都是8B=2单位
 */
func NewMemory(mem int) *Memory {
	var total = mem * 1024 * 1024 / 4
	var m = &Memory{Left: total}
	return m
}

type Memory struct {
	Left int
	Mtx  sync.Mutex
}

func (this *Memory) AskFor(count int) (success bool) {
	success = true
	this.Mtx.Lock()
	if this.Left >= count {
		this.Left -= count
	} else {
		success = false
	}
	this.Mtx.Unlock()
	return
}

func (this *Memory) GiveBack(count int) {
	this.Mtx.Lock()
	this.Left += count
	this.Mtx.Unlock()
}

//bottom
//
