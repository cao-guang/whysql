package main

import "strings"

var _singleGetTemplate = `
// NAME {{or .Comment "get data from rc"}} 
func (d *{{.StructName}}) NAME(c context.Context,h KEYSS{{.ExtraArgsType}}) (res VALUE, err error) {
	var bs []byte
	conn := d.redis.Conn(c)
	cacheKey := {{.KeyMethod}}({{.ExtraArgs}})
	cacheKey = library.XT_HCBS + ":" + cacheKey
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

//, "EX", {{.ExpireCode}} 去掉设置缓存过期时间变为永不过期
var _singleSetTemplate = `
// NAME {{or .Comment "Set data to rc"}} 
func (d *{{.StructName}}) NAME(c context.Context, h KEYSS, data VALUE {{.ExtraArgsType}}) (err error) {
	var bs []byte
	conn := d.redis.Conn(c)
	cacheKey := {{.KeyMethod}}({{.ExtraArgs}})
	cacheKey = library.XT_HCBS + ":" + cacheKey
	defer conn.Close()
	if bs, err = json.Marshal(data); err != nil {
		log.Error("json.Marshal(%+v) error(%v)", data, err)
		return
	}
	_, err =conn.Do("SET", cacheKey, bs)
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
	iter := 0
	var keys []string
	for {
		if arr, err := red.MultiBulk(conn.Do("SCAN", iter,"MATCH",library.XT_HCBS +"*"+key+"*")); err != nil {
			panic(err)
		} else {
			iter, _ = red.Int(arr[0], nil)
			key,_:=red.Strings(arr[1], nil)
			for _, value := range key {
				keys=append(keys,value)
			}
		}
		if iter == 0  {
			break
		}
	} 
	for i, _ := range keys {
		conn.Send("DEL", keys[i])
	} 
	if err!=nil{
		log.Error("NAME conn.Close error(%v)", err)
	}
	return
}
`


