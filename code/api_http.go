package main

import (
	`encoding/json`
	`fmt`
	`greentea/utils`
	`io/ioutil`
	"net/http"
	`strconv`
	`strings`
)

var (
	_handler_map = make(map[string]func(*http.Request, []byte) interface{}, 16)

	getKey func(string) uint32
)

func init() {
	_handler_map[`mset`] = iMSet
	_handler_map[`mget`] = iMGet
	_handler_map[`push`] = lMPush
	_handler_map[`pop`] = lPop
	_handler_map[`spush`] = sMPush
	_handler_map[`spop`] = sPop
	_handler_map[`zadd`] = zAdd
	_handler_map[`zpush`] = zMPush
	_handler_map[`zpop`] = zPop
	_handler_map[`switch`] = switchDb
	_handler_map[`ctrl`] = control
	http.HandleFunc(`/`, httpEntrance)

}

func httpInit() {
	if _support_str_key {
		getKey = getKeyStr
	} else {
		getKey = getKeyInt
	}
}

func getKeyStr(str string) uint32 {
	var key, err = strconv.ParseUint(str, 10, 32)
	if err != nil { //如果key是字符串，则hash到整数
		key = utils.BKDRHash(str)
	}
	return uint32(key)
}
func getKeyInt(str string) uint32 {
	var key, err = strconv.ParseUint(str, 10, 32)
	if err != nil {
		throw(err)
	}
	return uint32(key)
}

func getDb(r *http.Request) uint32 {
	var str_db = r.Form.Get(`db`)
	var val, err = strconv.ParseUint(str_db, 10, 32)
	if err == nil {
		return uint32(val)
	}
	return _default_db
}

func sendJson(w http.ResponseWriter, data interface{}) {
	var bytes, err = json.Marshal(data)
	throw(err, `result to json`)
	fmt.Fprint(w, `{"err_no":0,"result":`, string(bytes), `}`)
}

func sendError(w http.ResponseWriter, msg interface{}) {
	fmt.Fprint(w, `{"err_no":1,"err_msg":"`, msg, `"}`)
}

//没用上
type httpReq struct {
	w    http.ResponseWriter
	r    *http.Request
	body []byte
	at   *AllTable
	db   uint32
}

func httpEntrance(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if e := recover(); e != nil {
			sendError(w, e)
		}
	}()
	var uri = r.URL.RequestURI()[1:]
	var pos = strings.Index(uri, `?`)
	if pos < 0 {
		pos = len(uri)
	}
	var handler, ok = _handler_map[uri[:pos]]
	if ok {
		r.ParseForm()
		var body, err = ioutil.ReadAll(r.Body)
		throw(err)
		var ret = handler(r, body)
		sendJson(w, ret)
	} else {
		sendError(w, err_operation)
	}
}

func switchDb(r *http.Request, body []byte) interface{} {
	var db = getDb(r)
	if db > _db_number {
		panic(err_out_range)
	}
	_default_db = db
	return db
}

func control(r *http.Request, body []byte) interface{} {
	var db = getDb(r)
	var cmd = r.Form.Get(`cmd`)
	if cmd == `all` {
		var t = r.Form.Get(`t`)
		if t == `` {
			t = `i`
		}
		var limit uint32 = 10000
		switch t {
		case `i`:
			var ret = map[string]int64{}
			var tb = _all_tables[db-1].it
			var ln = uint32(len(tb.Keys))
			if limit > ln {
				limit = ln
			}
			var keys = tb.Keys[:limit]
			var vals = tb.Vals[:limit]
			for i, key := range keys {
				var str = strconv.Itoa(int(key))
				ret[str] = vals[i]
			}
			return ret
		case `l`, `s`:
			var ret = map[string][]int64{}
			var tb = _all_tables[db-1].GetLT(t)
			var ln = uint32(len(tb.Keys))
			if limit > ln {
				limit = ln
			}
			var keys = tb.Keys[:limit]
			var vals = tb.Vals[:limit]
			for i, key := range keys {
				var str = strconv.Itoa(int(key))
				ret[str] = vals[i].FetchAll()
			}
			return ret
		case `z`:
			var ret = map[string][]Unit{}
			var tb = _all_tables[db-1].zt
			var ln = uint32(len(tb.Keys))
			if limit > ln {
				limit = ln
			}
			var keys = tb.Keys[:limit]
			var vals = tb.Vals[:limit]
			for i, key := range keys {
				var str = strconv.Itoa(int(key))
				ret[str] = vals[i].FetchAll()
			}
			return ret
		default:
		}
	}
	return 1
}
