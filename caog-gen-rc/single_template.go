package main

import "strings"

var _singleGetTemplate = `
// NAME {{or .Comment "get data from rc"}} 
func (d *{{.StructName}}) NAME(c context.Context,h KEYSS{{.ExtraArgsType}}) (res VALUE, err error) {
	var bs []byte
	conn := d.redis.Conn(c)
	cacheKey := {{.KeyMethod}}({{.ExtraArgs}})
	cacheKey = library.XT_HCBS + ":EXKEY:" + cacheKey
	defer conn.Close()
	bs, err = redis.Bytes(conn.Do("GET", cacheKey))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
		}
		return
	}
	gzipData, err := library.GzipDecode(bs)
	res = new({{.OriginValueType}})
	if err = json.Unmarshal(gzipData, res); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", string(bs), err)
		return
	}
	return
}
`

// , "EX", {{.ExpireCode}} 去掉设置缓存过期时间变为永不过期
var _singleSetTemplate = `
// NAME {{or .Comment "Set data to rc"}} 
func (d *{{.StructName}}) NAME(c context.Context, h KEYSS, data VALUE {{.ExtraArgsType}}) (err error) {
	var bs []byte
	conn := d.redis.Conn(c)
	cacheKey := {{.KeyMethod}}({{.ExtraArgs}})
	//先对key进行切割给每一个关联的地方都添加进去详细的key
	arrkey := strings.Split(key, ":")
	cacheKey = library.XT_HCBS + ":EXKEY:" + cacheKey
	defer conn.Close()
	for _, v := range arrkey {
		if library.IsNum(v) || v[0:1] == "0" {
			continue
		}else{
			err = conn.Send("SADD", library.XT_HCBS+":DELKEY:SETKEY:"+v, cacheKey)
			if err != nil {
				log.Error("NAME conn.Send(SET, %s) error(%v)", cacheKey, err)
				return
			}
			if err = conn.Send("EXPIRE", library.XT_HCBS+":DELKEY:SETKEY:"+v, 86400); err != nil {
				log.Error("NAME conn.Send error(%v)", err)
				return
			} 
		} 
	}
	if bs, err = json.Marshal(data); err != nil {
		log.Error("json.Marshal(%+v) error(%v)", data, err)
		return
	}
	gzipData, err := library.GzipEncode(bs)
	err = conn.Send("SET", cacheKey, gzipData,"EX",86400)
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
	keys := library.XT_HCBS + ":DELKEY:SETKEY:" + key
	r2, err := conn.Do("SMEMBERS", keys)
	if err != nil {
		log.Error("NAME conn.Do(SGET, %s) error(%v)", keys, err)
		return
	}
	temp, ok := r2.([]interface{})
	if !ok {
		log.Error("NAME conn.Close error(%v)", err)
		return
	}
	for _, v := range temp {
		err := conn.Send("DEL", v)
		if err != nil {
			log.Error("NAME conn.Do(DEL, %s) error(%v)", keys, err)
			continue
		} else {
			err = conn.Send("SREM", keys, v)
			if err != nil {
				log.Error("NAME conn.Do(SREM, %s) error(%v)", keys, err)
				continue
			}
		}
	}
	return
}
`
