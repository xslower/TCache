package main

import (
	`errors`
	`flag`
	`github.com/astaxie/beego/config`
	`greentea/utils`
	`net/http`
	`strconv`
)

var (
	err_no_data   = errors.New(`No Data in here`)
	err_no_enough = errors.New(`No enough data`)
	err_out_range = errors.New(`Index is out of range`)
	err_format    = errors.New(`Not support format!`)
	err_operation = errors.New(`Not Support Operation!`)

	_max_list_length uint32 = 10000
	_max_memory      uint32 = 128
	_db_number       uint32 = 16
	_default_db      uint32 = 1
	_port            string = `:6789`
	_support_str_key bool   = true

	_all_tables = []*AllTable{}
)

func init() {
	conf_file := flag.String(`c`, `config.ini`, `-c /path/to/config.ini`)
	flag.Parse()
	cnf, err := config.NewConfig(`ini`, *conf_file)
	throw(err, `configure file error!`)
	global, err := cnf.GetSection(`global`)
	throw(err, `configure format error!`)
	utils.Init(global[`log_file`])
	str := global[`port`]
	if str != `` {
		_port = `:` + str
	}
	str = global[`support_string_key`]
	if str == `` || str == `0` || str == `false` {
		_support_str_key = false
	}
	val, err := strconv.ParseUint(global[`max_list_length`], 10, 32)
	if err == nil {
		_max_list_length = uint32(val)
	}
	val, err = strconv.ParseUint(global[`max_memory`], 10, 32)
	if err == nil {
		_max_memory = uint32(val)
	}
	val, err = strconv.ParseUint(global[`db_number`], 10, 32)
	if err == nil {
		_db_number = uint32(val)
	}
	httpInit()
}

func main() {

	initDb(_max_memory, _db_number)
	err := http.ListenAndServe(_port, nil)
	if err != nil {
		println(err.Error())
	}

}

//bottom
