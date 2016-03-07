/**
 * 所有的方法遍历两次map，是为了可以在string转int失败时抛出异常，
 * 而不会插入部分数据
 */

package main

import (
	`encoding/json`
	"net/http"
	`strconv`
)

func iIncr(r *http.Request, body []byte) interface{} {
	var ap = &ApiParam{}
	var err = json.Unmarshal(body, ap)
	throw(err)
	var db = getDb(r)
	var tb = _all_tables[db-1].it
	return tb.Increase(ap.Key, ap.Val)

}

func intMSet(r *http.Request, body []byte, t string) interface{} {
	var mapPost = map[string]int64{}
	var err = json.Unmarshal(body, &mapPost)
	var ln = len(mapPost)
	throw(err)
	var formated = make(map[uint32]int64, ln)
	var key uint32 = 0
	for str_key, val := range mapPost {
		key = getKey(str_key)
		formated[key] = val
	}
	var db = getDb(r)
	var atb = _all_tables[db-1]
	var tb ITbSet
	switch t {
	case `i`:
		tb = atb.it
	case `l`:
		tb = atb.lt
	case `s`:
		tb = atb.st
	}
	for k, v := range formated {
		tb.Set(k, v)
	}
	return ln
}

/**
 * int multi set
 */
func iMSet(r *http.Request, body []byte) interface{} {
	return intMSet(r, body, `i`)
}

/**
 * list multi set
 */
func lMPush(r *http.Request, body []byte) interface{} {
	return intMSet(r, body, `l`)
}

/**
 * sorted list multi set
 */
func sMPush(r *http.Request, body []byte) interface{} {
	return intMSet(r, body, `s`)
}

type KeyMark struct {
	Key uint32
	Unit
}

func zAdd(r *http.Request, body []byte) interface{} {
	var kmark = &KeyMark{}
	var err = json.Unmarshal(body, kmark)
	throw(err)
	var db = getDb(r)
	var tb = _all_tables[db-1].zt
	tb.Set(kmark.Key, kmark.Unit)
	return 1
}

/**
 * sorted set multi set
 */
func zMPush(r *http.Request, body []byte) interface{} {
	var kmarks = []KeyMark{}
	var err = json.Unmarshal(body, &kmarks)
	throw(err)
	var db = getDb(r)
	var tb = _all_tables[db-1].zt
	_ = "breakpoint"
	for _, mk := range kmarks {
		tb.Set(mk.Key, mk.Unit)
	}
	return len(kmarks)
}

/**
 * sorted list value-keys multi set.
 * format: val:[key1,key2,key3]
 */
func sVKPush(r *http.Request, body []byte) interface{} {
	var mapPost = map[string][]uint32{}
	var err = json.Unmarshal(body, &mapPost)
	throw(err)
	var formated = map[uint32]int64{}
	var val int64 = 0
	for str_val, slc_keys := range mapPost {
		val, err = strconv.ParseInt(str_val, 10, 64)
		throw(err)
		for _, key := range slc_keys {
			formated[key] = val
		}
	}
	var db = getDb(r)
	var tb = _all_tables[db-1].st
	for k, v := range formated {
		tb.Set(k, v)
	}
	return len(formated)
}

/**
 * bottom
 */
