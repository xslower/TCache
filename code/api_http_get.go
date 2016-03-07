package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func iMGet(r *http.Request, body []byte) interface{} {
	var arry = []uint32{}
	var err = json.Unmarshal(body, &arry)
	throw(err)
	var db = getDb(r)
	var withKey = r.Form.Get(`withkey`)
	var it = _all_tables[db-1].it
	if withKey != `` {
		var ret = make(map[string]int64, len(arry))
		for _, key := range arry {
			var val = it.Get(key)
			var str_key = strconv.FormatUint(uint64(key), 10)
			ret[str_key] = val
		}
		return ret
	} else {
		var ret = make([]int64, len(arry))
		for i, key := range arry {
			var val = it.Get(key)
			ret[i] = val
		}
		return ret
	}
}

type ApiParam struct {
	Key   uint32
	Val   int64
	Score float32
	Right int8
}

func slPop(r *http.Request, body []byte, t string) interface{} {
	var pp = &ApiParam{}
	var err = json.Unmarshal(body, pp)
	throw(err)
	var db = getDb(r)
	var tb = _all_tables[db-1].GetLT(t)
	var list = tb.Get(pp.Key)
	var num = uint32(pp.Val)
	var ret []int64
	if pp.Right > 0 {
		ret = list.PopRightMulti(num)
	} else {
		ret = list.PopLeftMulti(num)
	}
	return ret
}

func lPop(r *http.Request, body []byte) interface{} {
	return slPop(r, body, `l`)
}

func sPop(r *http.Request, body []byte) interface{} {
	return slPop(r, body, `s`)
}

func zPop(r *http.Request, body []byte) interface{} {
	var pp = &ApiParam{}
	var err = json.Unmarshal(body, pp)
	throw(err)
	var db = getDb(r)
	var tb = _all_tables[db-1].zt
	var set = tb.Get(pp.Key)
	var num = uint32(pp.Val)
	var ret []Unit
	if pp.Right > 0 {
		ret = set.PopRightMulti(num)
	} else {
		ret = set.PopLeftMulti(num)
	}
	return ret
}

type FetchParam struct {
	Key    uint32
	ISmall int64
	IBig   int64
	FSmall float32
	FBig   float32
	Limit  int
	Odr    Order
}

func sFetchBetween(r *http.Request, body []byte) interface{} {
	var fp = &FetchParam{}
	var err = json.Unmarshal(body, fp)
	throw(err)
	var db = getDb(r)
	var tb = _all_tables[db-1].st
	var list = tb.Get(fp.Key)
	var ret = list.FetchBetweenLimit(fp.ISmall, fp.IBig, fp.Limit, fp.Odr)
	return ret
}

func zFetchBetween(r *http.Request, body []byte) interface{} {
	var fp = &FetchParam{}
	var err = json.Unmarshal(body, fp)
	throw(err)
	var db = getDb(r)
	var tb = _all_tables[db-1].zt
	var list = tb.Get(fp.Key)
	var ret = list.FetchBetweenLimit(fp.FSmall, fp.FBig, fp.Limit, fp.Odr)
	return ret
}

func aMDel(r *http.Request, body []byte, t string) interface{} {
	var slcKeys = []uint32{}
	var err = json.Unmarshal(body, &slcKeys)
	throw(err)
	var db = getDb(r)
	var tb = _all_tables[db-1].GetIDel(t)
	for _, key := range slcKeys {
		tb.Del(key)
	}
	return len(slcKeys)
}

func iMDel(r *http.Request, body []byte) interface{} {
	return aMDel(r, body, `i`)
}

func lMDel(r *http.Request, body []byte) interface{} {
	return aMDel(r, body, `l`)
}

func sMDel(r *http.Request, body []byte) interface{} {
	return aMDel(r, body, `s`)
}

func zMDel(r *http.Request, body []byte) interface{} {
	return aMDel(r, body, `z`)
}
