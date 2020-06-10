package main

import "strings"

var _singleGetTemplate = `
// NAME {{or .Comment "get data from rc"}} 
func (d *{{.StructName}}) NAME(c context.Context,h KEYSS{{.ExtraArgsType}}) (res VALUE, err error) {
	var bs []byte
	conn := d.redis.Conn(c)
	cacheKey := {{.KeyMethod}}({{.ExtraArgs}})
	defer conn.Close()
	bs, err = redis.Bytes(conn.Do("GET", cacheKey))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
		}
		return
	}
	res = new({{.OriginValueType}})
	if err = json.Unmarshal(bs, res); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", string(bs), err)
		return
	}
	return
}
`

var _singleSetTemplate = `
// NAME {{or .Comment "Set data to rc"}} 
func (d *{{.StructName}}) NAME(c context.Context, h KEYSS, data VALUE {{.ExtraArgsType}}) (err error) {
	var bs []byte
	conn := d.redis.Conn(c)
	cacheKey := {{.KeyMethod}}({{.ExtraArgs}})
	defer conn.Close()
	if bs, err = json.Marshal(data); err != nil {
		log.Error("json.Marshal(%+v) error(%v)", data, err)
		return
	}
	_, err =conn.Do("SET", cacheKey, bs, "EX", {{.ExpireCode}})
	if err != nil {
		log.Error("NAME conn.Send(SET, %s) error(%v)", cacheKey, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("NAME conn.Flush error(%v)", err)
		return
	}
	return
}
`

var _singleAddTemplate = strings.Replace(_singleSetTemplate, "Set", "Add", -1)
var _singleReplaceTemplate = strings.Replace(_singleSetTemplate, "Set", "Replace", -1)

var _singleDelTemplate = `
// NAME {{or .Comment "delete data from rc"}} 
func (d *{{.StructName}}) NAME(c context.Context, key KEYSS {{.ExtraArgsType}}) (err error) {
	r := d.redis
	conn := r.Conn(c)
	defer  conn.Close()
	val, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	conn.Send("MULTI")
	for i, _ := range val {
		conn.Send("DEL", val[i])
	}
	_, err = redis.Values(conn.Do("EXEC"))
	if err!=nil{
		log.Error("NAME conn.Close error(%v)", err)
	}
	return
}
`


