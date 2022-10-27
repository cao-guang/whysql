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
	//先对key进行切割给每一个关联的地方都添加进去详细的key
  	arrkey := strings.Split(key,":")
	cacheKey = library.XT_HCBS + ":" + cacheKey
	defer conn.Close()
	for _,v:=range arrkey{
		_, err = conn.Do("SADD", library.XT_HCBS + ":DELKEY:SETKEY:" + v,cacheKey)
		if err != nil {
			  log.Error("NAME conn.Send(SET, %s) error(%v)", cacheKey, err)
			  return
		}
  	}
  	if bs, err = json.Marshal(data); err != nil {
		log.Error("json.Marshal(%+v) error(%v)", data, err)
		return
  	}
  	_, err = conn.Do("SET", cacheKey, bs)
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
		log.Error("DelCacheHT_GET_YH_YH conn.Do(SGET, %s) error(%v)", keys, err)
		return
  	}
  	temp,ok:=r2.([]interface{})
 	 if !ok{
		log.Error("DelCacheHT_GET_YH_YH conn.Close error(%v)", err)
		return
 	 }
  	for _,v:=range temp{
		err := conn.Send("DEL", v)
		if err != nil {
			log.Error("DelCacheHT_GET_YH_YH conn.Do(DEL, %s) error(%v)", keys, err)
			continue
		}else{
		   err = conn.Send("SREM", keys,v)
		   if err != nil {
				log.Error("DelCacheHT_GET_YH_YH conn.Do(SREM, %s) error(%v)", keys, err)
				continue
		   }
		}
  	}
  	return
}
`


