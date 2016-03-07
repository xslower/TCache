package main

import (
	"greentea/utils"
	"time"
)

func getNextStep(cur int) int {
	if cur < 100 {
		return 100
	} else if cur < 1000 {
		return cur
	}
	cur = cur * 30 / 100
	return cur
}

func throw(err error, msg ...string) {
	utils.Throw(err, msg...)
}

func check(err error, msg ...interface{}) {
	utils.Check(err, msg...)
}

func logit(data ...interface{}) {
	utils.Logit(data...)
}

func echo(i ...interface{}) {
	utils.Echo(i...)
}

func echoStrSlice(strs ...[]string) {
	utils.EchoStrSlice(strs...)
}

func echoBytes(args interface{}) {
	utils.EchoBytes(args)
}

func getTimestamp() uint32 {
	return uint32(time.Now().Unix())
}
