package whysql

import (

	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gohouse/gorose/v2"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"


)

const (
	STUDENT   = "735800" //职业农民
	TEACHER   = "735900" //注册用户
	PARENT    = "735700" //农技人员
	ZONE_JSZJ = "735600" //技术专家
	Prev      = "pre_"

)

func buildPadStr(str string, padStr string, padLen int, padLeft bool, padRight bool) string {

	// When padded length is less then the current string size
	if padLen < utf8.RuneCountInString(str) {
		return str
	}

	padLen -= utf8.RuneCountInString(str)

	targetLen := padLen

	targetLenLeft := targetLen
	targetLenRight := targetLen
	if padLeft && padRight {
		targetLenLeft = padLen / 2
		targetLenRight = padLen - targetLenLeft
	}

	strToRepeatLen := utf8.RuneCountInString(padStr)

	repeatTimes := int(math.Ceil(float64(targetLen) / float64(strToRepeatLen)))
	repeatedString := strings.Repeat(padStr, repeatTimes)

	leftSide := ""
	if padLeft {
		leftSide = repeatedString[0:targetLenLeft]
	}

	rightSide := ""
	if padRight {
		rightSide = repeatedString[0:targetLenRight]
	}

	return leftSide + str + rightSide
}

// PadLeft pad left side of string if size of string is less then indicated pad length
func PadLeft(str string, padStr string, padLen int) string {
	return buildPadStr(str, padStr, padLen, true, false)
}

// PadRight pad right side of string if size of string is less then indicated pad length
func PadRight(str string, padStr string, padLen int) string {
	return buildPadStr(str, padStr, padLen, false, true)
}

//用户装载handle，调用dalsql方法
var EventMap = map[string]interface{}{}

func RegisterEvent(key string, fn interface{}) {
	EventMap[key] = fn
}

func Struct2Map(obj interface{}) map[string]interface{} {
	obj_v := reflect.ValueOf(obj)
	v := obj_v.Elem()
	typeOfType := v.Type()
	var data = make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		data[typeOfType.Field(i).Name] = field.Interface()
	}
	return data
}

//排除list的
func Struct2Map2(obj interface{}, list []string) map[string]interface{} {
	obj_v := reflect.ValueOf(obj)
	v := obj_v.Elem()
	typeOfType := v.Type()

	var data = make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		name := typeOfType.Field(i).Name
		val := field.Interface()
		if len(list) > 0 {
			findit := false
			for _, s := range list {
				if s == name {
					findit = true
					break
				}
			}
			if findit == false {
				data[name] = val
			}
		} else {
			data[name] = val
		}
	}
	return data
}

//返回list存在的
func Struct2Map3(obj interface{}, list []string) map[string]interface{} {
	obj_v := reflect.ValueOf(obj)
	v := obj_v.Elem()
	typeOfType := v.Type()

	var data = make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		name := typeOfType.Field(i).Name
		val := field.Interface()
		if len(list) > 0 {
			findit := false
			for _, s := range list {
				if s == name {
					findit = true
					break
				}
			}
			if findit == true {
				data[name] = val
			}
		} else {
			data[name] = val
		}
	}
	return data
}

//返回map在对象里面的map
func RetrunStructMap(obj interface{}, m map[string]interface{}) map[string]interface{} {
	obj_v := reflect.ValueOf(obj)
	v := obj_v.Elem()
	typeOfType := v.Type()
	var data = make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		name := typeOfType.Field(i).Name
		if len(m) > 0 {
			for k, v := range m {
				if strings.ToUpper(k) == name {
					data[name] = v
				}
			}
		}
	}
	return data
}


func TimeNowString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
func TimeMinString() string {
	return "0001/01/01"
}
func GetRJXDM(xk_id string) string {
	var rjxdm string
	if xk_id != "" && len(xk_id) == 6 {
		_xd := xk_id[3:4]
		switch _xd {
		case "1": //小学
			rjxdm = "178100"
			break
		case "2": //初中
			rjxdm = "178200"
			break
		case "4": //高中
			rjxdm = "178400"
			break
		}
	}
	return rjxdm
}

type IDataHelper interface {
	//新增的方法
	Add(interface{}) (string, error)
	//添加或者修改的方法，[]string,表示要修改要排除的字段
	AddOrUpdate(interface{}, []string, []string) (string, error)
	//根据主键去修改具体的值,
	UpdateField(interface{}, []string) (string, error)
	//根据where条件修改具体的值,这种是list方式
	UpdateFieldBySQL(param interface{}, list []string, where string) (string, error)
	//根据where条件修改具体的值,这种是map方式
	UpdateFieldMapByKey(param interface{}, m map[string]interface{}) (string, error)
	//根据where条件修改具体的值,这种是map方式
	UpdateFieldMap(param interface{}, m map[string]interface{}, where string) (string, error)
	//事务删除
	DeleteTran(map[string]string, string) (string, error)
	//获取对象
	GetModelInfo_MySQL(obj interface{}) error
	//获取系统配置
	GetXT_XTPZ(string) (string, error)
}
type DataHelper struct {
	DB 	gorose.IOrm

}

func CreateDataHelper(session gorose.IOrm) *DataHelper {
	d := &DataHelper{DB: session}
	return d
}

// 实现DataHelper接口的Add方法
func (d *DataHelper) Add(param interface{}) (string, error) {
	var strObj string
	data := Struct2Map(param)
	RemoveEmptyMap(data)
	count, err := d.DB.Table(param).Data(data).Insert()
	PrintSQL(d.DB.LastSql)
	if err != nil {
		if count > 0 {
			strObj = BackOk("新增成功")
		} else {
			strObj = BackErr("新增失败", err)
		}
	} else {
		strObj = BackOk("新增成功")
	}
	return strObj, err
}

// 添加或者修改， updatelist 表 要排除某些要修改的列,addlist 表 要排除某些要新增的列
func (d *DataHelper) AddOrUpdate(param interface{}, updatelist []string, addlist []string) (string, error) {
	var strObj string
	var err error
	var isadd bool = true

	//dstVal := reflect.ValueOf(param)
	//sliceVal := reflect.Indirect(dstVal)
	//tablename := sliceVal.Type().Name()
	obj_v := reflect.ValueOf(param)
	v := obj_v.Elem()
	typeOfType := v.Type()
	tablename := typeOfType.Name()
	key := typeOfType.Field(0).Name
	value := v.Field(0).Interface()
	var count int64
	var msg string = "新增"
	if value == "" {
		data := Struct2Map2(param, addlist)
		RemoveEmptyMap(data)
		count, err = d.DB.Table(tablename).Data(data).Insert()
	} else {
		m, _ := d.DB.Table(tablename).Where(key, value).First()
		if m != nil {
			data := Struct2Map2(param, updatelist) //获取需要修改的map对象，最后跟m合并下，变化的才去修改
			lastdata := UnionDiffMap(m, data)
			_, ok := lastdata[key]
			if ok {
				delete(lastdata, key)
			}
			msg = "修改"
			count, err = d.DB.Table(tablename).Data(lastdata).Where(key, value).Update()
			isadd = false
		} else {
			data := Struct2Map2(param, addlist)
			RemoveEmptyMap(data)
			count, err = d.DB.Table(tablename).Data(data).Insert()
		}
	}
	println(d.DB.LastSql)
	if isadd && d.DB.LastInsertId() > 0 { //给主键对象赋值
		v.Field(0).SetString(strconv.FormatInt(d.DB.LastInsertId(), 10))
	}
	if err != nil {
		if count > 0 {
			strObj = BackOkdata(msg+"成功", strconv.FormatInt(d.DB.LastInsertId(), 10))
		} else {
			strObj = BackErr(msg+"失败", err)
		}
	} else {
		strObj = BackOkdata(msg+"成功", strconv.FormatInt(d.DB.LastInsertId(), 10))
	}
	return strObj, err
}

// 添加或者修改， updatelist 表 要排除某些要修改的列,addlist 表 要排除某些要新增的列,upmap要修改某些的map
func (d *DataHelper) AddOrUpdateMap(param interface{}, updatelist []string, addlist []string, upmap map[string]interface{}) (string, error) {
	var strObj string
	var err error
	var isadd bool = true

	//dstVal := reflect.ValueOf(param)
	//sliceVal := reflect.Indirect(dstVal)
	//tablename := sliceVal.Type().Name()
	obj_v := reflect.ValueOf(param)
	v := obj_v.Elem()
	typeOfType := v.Type()
	tablename := typeOfType.Name()
	key := typeOfType.Field(0).Name
	value := v.Field(0).Interface()
	var count int64
	var msg string = "新增"
	if value == "" {
		data := Struct2Map2(param, addlist)
		RemoveEmptyMap(data)
		count, err = d.DB.Table(tablename).Data(data).Insert()
	} else {
		arr := strings.Split(value.(string), ",")

		if len(arr) > 1 {
			msg = "修改"
			data := RetrunStructMap(param, upmap)
			if data != nil {
				delete(data, key)
			}
			var arrs []interface{}
			for _,v:=range arr{
				arrs = append(arrs, v)
			}
			count, err = d.DB.Table(tablename).Data(data).WhereIn(key, arrs).Update()
		} else {
			m, _ := d.DB.Table(tablename).Where(key, value).First()
			if m != nil {
				data := Struct2Map2(param, updatelist) //获取需要修改的map对象，最后跟m合并下，变化的才去修改
				lastdata := UnionDiffMap(m, data)      //得到差异的map
				if upmap != nil && len(upmap) > 0 {
					for k, _ := range upmap { //传递过来的需要修改的map
						_, ok := lastdata[k]
						if ok == false {
							delete(lastdata, k)
						}
					}
					for k, _ := range lastdata {
						_, ok := upmap[k]
						if ok == false {
							delete(lastdata, k)
						}
					}
				}
				//lastdata = UnionDiffMap(lastdata, upmap)
				_, ok := lastdata[key]
				if ok {
					delete(lastdata, key)
				}
				msg = "修改"
				if len(lastdata) > 0 {
					count, err = d.DB.Table(tablename).Data(lastdata).Where(key, value).Update()
				}
				isadd = false
			} else {
				data := Struct2Map2(param, addlist)
				RemoveEmptyMap(data)
				count, err = d.DB.Table(tablename).Data(data).Insert()
			}
		}
	}
	println(d.DB.LastSql)
	if isadd && d.DB.LastInsertId() > 0 { //给主键对象赋值
		v.Field(0).SetString(strconv.FormatInt(d.DB.LastInsertId(), 10))
	}
	if err != nil {
		if count > 0 {
			strObj = BackOkdata(msg+"成功", strconv.FormatInt(d.DB.LastInsertId(), 10))
		} else {
			strObj = BackErr(msg+"失败", err)
		}
	} else {
		strObj = BackOkdata(msg+"成功", strconv.FormatInt(d.DB.LastInsertId(), 10))
	}
	return strObj, err
}


//list是必填参数
func (d *DataHelper) UpdateField(param interface{}, list []string) (string, error) {
	var err error
	obj_v := reflect.ValueOf(param)
	v := obj_v.Elem()
	typeOfType := v.Type()
	tablename := typeOfType.Name()
	key := typeOfType.Field(0).Name
	value := v.Field(0).Interface()
	//var count int64
	var msg string = "修改"
	data := Struct2Map3(param, list)
	arr := strings.Split(value.(string), ",")
	if len(arr) > 1 {
		var arrs []interface{}
		arrs = append(arrs, arr)
		_, err = d.DB.Table(tablename).Data(data).WhereIn(key, arrs).Update()
	} else {
		_, err = d.DB.Table(tablename).Data(data).Where(key, value).Update()
	}
	PrintSQL(d.DB.LastSql)
	return ActionBack(msg, err)
}
func (d *DataHelper) UpdateFieldBySQL(param interface{}, list []string, where string) (string, error) {
	var err error
	sql, _ := GetUpdateFieldBySQL(param, list, nil, where)
	_, err = d.DB.Execute(sql)
	return ActionBack("修改", err)
}
func (d *DataHelper) UpdateFieldMapByKey(param interface{}, data map[string]interface{}) (string, error) {
	obj_v := reflect.ValueOf(param)
	v := obj_v.Elem()
	typeOfType := v.Type()
	key := typeOfType.Field(0).Name
	value := v.Field(0).Interface()
	where := key + "= '" + value.(string) + "'"
	return d.UpdateFieldMap(param, data, where)
}

func (d *DataHelper) UpdateFieldMap(param interface{}, data map[string]interface{}, where string) (string, error) {
	sql, _ := GetUpdateFieldBySQLMap(param, data, nil, where)
	_, err := d.DB.Execute(sql)
	return ActionBack("修改", err)
}

func ActionBack(msg string, err error) (string, error) {
	var strObj string
	if err != nil {
		strObj = BackErr(msg+"失败", err)
	} else {
		strObj = BackOk(msg + "成功")
	}
	return strObj, err
}
func (d *DataHelper) DeleteTran(m map[string]string, val string) (string, error) {
	var msg string
	msg = "操作成功"
	errn := d.DB.Transaction(func(db gorose.IOrm) error {
		for k, v := range m {
			sql := fmt.Sprintf("DELETE FROM %s WHERE  %s in (%s) ", k, v, FQ_instr(val))
			PrintSQL(sql)
			_, err2 := d.DB.Execute(sql)
			if err2 != nil {
				msg = "删除失败"
				return err2
			}
		}
		return nil
	})
	if errn != nil {
		return BackErr(msg, errn), errn
	} else {
		return BackOk(msg), nil
	}
}
func (d *DataHelper) GetModelInfo_MySQL(obj interface{}) error {
	obj_v := reflect.ValueOf(obj)
	v := obj_v.Elem()
	typeOfType := v.Type()
	tablename := typeOfType.Name()
	key := typeOfType.Field(0).Name
	value := v.Field(0).Interface()
	sql := "SELECT * FROM " + tablename + " WHERE " + key + "='" + value.(string) + "'"
	arr, err := d.DB.Query(sql)
	if err != nil {
		PrintSQL("GetModelInfo_MySQL", err)
		return err
	}
	if len(arr) > 0 {
		list := arr[0]
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			name := typeOfType.Field(i).Name
			if len(list) > 0 {
				for s, val := range list {
					if s == name && val != nil {
						var v interface{}
						switch val.(type) {
						case int64:
							v = fmt.Sprintf("%d", val)
						case float64:
							v, _ = strconv.ParseFloat(fmt.Sprintf("%f", val), 32)
						default:
							v = val
						}
						field.Set(reflect.ValueOf(v.(string)))
					}
				}

			}
		}
	}
	return nil
}

//返回m2存在m1的变化的数据,m2在m1里面的交集
func UnionDiffMap(m1 map[string]interface{}, m2 map[string]interface{}) map[string]interface{} {
	var m3 = make(map[string]interface{})
	if m1 != nil && m2 != nil {
		for k, v := range m1 {
			val, ok := m2[k]
			if ok {
				if val == "" && m1[k] == nil {
					continue
				} else if v != val { //如果val==""这种情况，界面上可能清空了，也必须修改为空
					m3[k] = val
				}
			}
		}
	}
	return m3
}

//删除空项
func RemoveEmptyMap(m map[string]interface{}) {
	if m != nil {
		for k, v := range m {
			if v == nil || v == "" {
				delete(m, k)
			}
		}
	}
}
func Map_GetParamsMap(h string, model interface{}) (string, map[string]interface{}, error) {
	m := make(map[string]interface{})
	m2 := make(map[string]interface{})
	if err := json.Unmarshal([]byte(h), &m); err == nil {
		if err := json.Unmarshal([]byte(h), &model); err != nil {
			return BackErr("序列化对象失败", nil), m2, errors.New("序列化对象失败")
		}
		//pars, _ok := m["parames"]
		//if _ok {
			str := h
			if err := json.Unmarshal([]byte(str), &m2); err == nil {
				for k, _ := range m2 {
					if strings.ToUpper(k) != k {
						m2[strings.ToUpper(k)] = m2[k] //把小写的赋值给大写，然后小写的
						delete(m2, k)
					}
				}
				return "", m2, nil
			}else{
				return "", m2, nil
			}
		//} else {
		//	return "", m2, nil
		//}
	}
	return BackErr("序列化对象失败", nil), m2, errors.New("序列化对象失败")
}

//param对象是有值的结构体对象
func String_GetParames(param interface{}, parames string) string {
	js, _ := json.Marshal(param)
	b := make(map[string]interface{})
	json.Unmarshal([]byte(js), &b)
	b["parames"] = string(parames)
	str, _ := json.Marshal(b)
	c := string(str)
	return c
}

//反射调用方法
func Call(m map[string]interface{}, name string, params ...interface{}) ([]reflect.Value, error) {
	f := reflect.ValueOf(m[name])
	if len(params) != f.Type().NumIn() {
		return nil, errors.New("参数不匹配!")
	}
	in := make([]reflect.Value, len(params))
	for k, v := range params {
		in[k] = reflect.ValueOf(v)
	}
	return f.Call(in), nil
}

//主要用于In查询，整形变字符串
func FQ_instr(str string) string {
	var res = str
	if str != "" {
		res = "'" + strings.Replace(str, ",", "','", -1) + "'"
	}
	return res
}

func PrintSQL(a ...interface{}) {
	//log.Println(a)
}

//替换所有位数的0
func ReplaceLast00(str string) string {
	pat := "[0]*$" //正则
	reg, _ := regexp.Compile(pat)
	return reg.ReplaceAllString(str, "")
}
func ReplaceLast00_DWH(QY_DM string) string {
	if QY_DM != "" {
		QY_DM = PadRight(QY_DM, "0", 12)
	}
	var qydm = ReplaceLast00(QY_DM)
	var parentid string = qydm
	var sXzqydmValue1 = ""
	lenparentid := len(parentid)
	if lenparentid == 0 {
		sXzqydmValue1 = ""
	} else if lenparentid >= 1 && lenparentid <= 2 { //省
		sXzqydmValue1 = QY_DM[0:2] + ""
	} else if lenparentid >= 3 && lenparentid <= 4 { //市
		sXzqydmValue1 = QY_DM[0:4] + ""
	} else if lenparentid >= 5 && lenparentid <= 6 { //县
		sXzqydmValue1 = QY_DM[0:6] + ""
	} else if lenparentid >= 7 && lenparentid <= 9 { //乡镇
		sXzqydmValue1 = QY_DM[0:9] + ""
	} else { //村社
		sXzqydmValue1 = QY_DM
	}
	return sXzqydmValue1
}
func GetQY_DM(QY_DM string) string {
	var qydm = ReplaceLast00(QY_DM)
	var parentid string = qydm
	var sXzqydmValue1 = ""
	lenparentid := len(parentid)
	if lenparentid == 0 {
		sXzqydmValue1 = "____________"
	} else if lenparentid >= 1 && lenparentid <= 2 { //省查询市
		sXzqydmValue1 = QY_DM[0:2] + "__________"
	} else if lenparentid >= 3 && lenparentid <= 4 { //市查县
		sXzqydmValue1 = QY_DM[0:4] + "________"
	} else if lenparentid >= 5 && lenparentid <= 6 {
		sXzqydmValue1 = QY_DM[0:6] + "______"
	} else if lenparentid >= 7 && lenparentid <= 9 {
		sXzqydmValue1 = QY_DM[0:9] + "___"
	} else {
		sXzqydmValue1 = QY_DM
	}
	return sXzqydmValue1
}
func GetSFZHInfo(SFZH string) (string, string, error) {
	var xb = "30701"
	var Birthday string = ""
	var strSex string = ""
	if SFZH == "" {
		return xb, "", nil
	} else {
		if len(SFZH) == 18 {
			Birthday = SFZH[6:10] + "/" + SFZH[10:12] + "/" + SFZH[12:14]
			strSex = SFZH[14:17]
		}
		if len(SFZH) == 15 {
			Birthday = "19" + SFZH[6:8] + "/" + SFZH[8:10] + "/" + SFZH[10:12]
			strSex = SFZH[12:15]
		}
		istrSex, err := strconv.Atoi(strSex)
		if err == nil {
			if istrSex%2 == 0 {
				xb = "30702"
			} else {
				xb = "30701"
			}
			return xb, Birthday, nil
		} else {
			return "", "", err
		}
	}
}

//得到SQL语句，List 表示要排除的字段，m表示需要格式化的对象
func GetInsertSql(obj interface{}, list []string, m map[string]interface{}) (string, error) {
	obj_v := reflect.ValueOf(obj)
	v := obj_v.Elem()
	typeOfType := v.Type()
	tablename := typeOfType.Name()
	var fileds []string
	var values []string
	//var data = make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		name := typeOfType.Field(i).Name
		val := field.Interface()
		if val.(string) != "" {
			var str string = "'" + val.(string) + "'"
			if m != nil && m[name] != nil {
				var a = new(Formatter)
				if m[name] == "date" {
					str = a.date(val.(string))
				}
			}
			if len(list) > 0 {
				findit := false
				for _, s := range list {
					if s == name {
						findit = true
						break
					}
				}
				if findit == false {
					fileds = append(fileds, name)
					values = append(values, str)

				}
			} else {
				fileds = append(fileds, name)
				values = append(values, str)
			}
		}

	}
	sql := " INSERT INTO " + tablename + " (" + strings.Join(fileds, ",") + ") VALUES (" + strings.Join(values, ",") + ")"
	PrintSQL(sql)
	return sql, nil
}

//获取要修改的SQL语句，通用，注意param 对象,list返回存在需要修改的, m主要用于格式化函数，
func GetUpdateFieldBySQL(param interface{}, list []string, m map[string]interface{}, where string) (string, error) {
	data := Struct2Map3(param, list)
	return GetUpdateFieldBySQLMap(param, data, m, where)
}

//获取要修改的SQL语句，通用，注意param 对象,m1返回存在需要修改的, m主要用于格式化函数，
func GetUpdateFieldBySQLMap(param interface{}, m1 map[string]interface{}, m map[string]interface{}, where string) (string, error) {
	var err error
	tname, ok := param.(string)
	tablename := ""
	fieldkey := ""
	if ok == true {
		tablename = tname
	} else {
		obj_v := reflect.ValueOf(param)
		v := obj_v.Elem()
		typeOfType := v.Type()
		tablename = typeOfType.Name()
		fieldkey = typeOfType.Field(0).Name //主键不允许修改
	}
	data := m1
	var arr []string
	var a = new(Formatter)
	for key, v := range data {
		var str string = "'" + v.(string) + "'"
		if m != nil && m[key] != nil {
			if m[key] == "date" {
				str = a.date(v.(string))
			}
		}
		if key != fieldkey {
			arr = append(arr, "t."+key+"="+str+"")
		}
	}

	if where == "" {
		where = " 1=2 "
	}
	s := strings.Join(arr, ",")
	if s != "" {
		var sql = " UPDATE " + tablename + " t SET " + s + " WHERE " + where
		PrintSQL(sql)
		return sql, err
	} else {
		return "", nil
	}

}

type Formatter struct {
}

func (c *Formatter) date(val string) string {
	return "to_date('" + val + "','yyyy-MM-dd HH24:mi:ss')"
}

//此类只是本go 使用
type BasePage struct {
	Page  int    `json:"page"`
	Rows  int    `json:"rows"`
	Sort  string `json:"sort"`
	Order string `json:"order"`
}

func (c *BasePage) Init() {
	if c.Page == 0 {
		c.Page = 1
	}
	if c.Rows == 0 {
		c.Rows = 20
	}
}

type MyError struct {
	error
}
type TryCatch struct {
	errChan      chan interface{}
	catches      map[reflect.Type]func(err error)
	defaultCatch func(err error)
}

func (t TryCatch) Try(block func()) TryCatch {
	t.errChan = make(chan interface{})
	t.catches = map[reflect.Type]func(err error){}
	t.defaultCatch = func(err error) {}
	go func() {
		defer func() {
			t.errChan <- recover()
		}()
		block()
	}()
	return t
}

func (t TryCatch) CatchAll(block func(err error)) TryCatch {
	t.defaultCatch = block
	return t
}

func (t TryCatch) Catch(e error, block func(err error)) TryCatch {
	errorType := reflect.TypeOf(e)
	t.catches[errorType] = block
	return t
}

func (t TryCatch) Finally(block func()) TryCatch {
	err := <-t.errChan
	if err != nil {
		catch := t.catches[reflect.TypeOf(err)]
		if catch != nil {
			catch(err.(error))
		} else {
			t.defaultCatch(err.(error))
		}
	}
	block()
	return t
}

func ThrowsPanic(f func()) (b bool) {
	defer func() {
		if x := recover(); x != nil {
			log.Println(x)
			b = true
		}
	}()
	f() //执行函数f，如果f中出现了panic，那么就可以恢复回来
	return b
}

type GoExchangeData struct {
	Exchangename string
	Routekey     string
	Msgbody      string
}

//请求goapi方法
func PostGoApi(godns string, method string, parames string) (string, error) {
	var objstr string
	csz := "parames=" + parames
	var url = strings.TrimSuffix(godns, "/") + "/" + strings.TrimPrefix(method, "/")

	//resp, err := http.Post(url,
	//	"application/x-www-form-urlencoded",
	//	strings.NewReader(csz))
	//if err != nil {
	//	log.Println(err)
	//}

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, strings.NewReader(csz))
	if err != nil {
		log.Println(err)
	}
	req.Header.Add("User-Agent", "golang post")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	objstr = string(body)
	fmt.Println("PostGoApi返回结果：" + objstr)
	return objstr, err
}
func GetCount(db gorose.IOrm, sql string) (m int64, errs error) {
	var err error

	pat := `(?i:order)[\s|\S]+(?i:by)[\s|\S]+$` //正则
	reg, _ := regexp.Compile(pat)
	sql = reg.ReplaceAllString(sql, "")
	sql = fmt.Sprintf("SELECT count(*) TOTAL FROM (%s) A", sql)
	resultSlice, err := db.Query(sql)
	if err != nil {
		return 0, err
	} else {
		if len(resultSlice) > 0 {
			a := fmt.Sprintf("%d", resultSlice[0]["TOTAL"])
			b, _ := strconv.ParseInt(a, 10, 64)
			return b, nil

		} else {
			return 0, nil
		}
	}
}
func GetPageRows(db gorose.IOrm, sql string, page int, rows int) (result []gorose.Data, errs error) {
	start := (page - 1) * rows
	var fenye = fmt.Sprintf(" limit %d,%d ", start, rows)
	sSql := sql + fenye
	PrintSQL("GetPageRows===>", sSql)
	return db.Query(sSql)
}
func QueryList(db gorose.IOrm, sql string, page int, rows int) (total string,result []gorose.Data, err error) {
	count, err := GetCount(db, sql)
	a01, err := GetPageRows(db, sql, page, rows)
	result=a01
	total = strconv.FormatInt(count,10)
	return total,result,err
}
//base64解密算法
func Base64DecodeString(str string) string {
	decoded, _ := base64.StdEncoding.DecodeString(str)
	str = string(decoded)
	return str
}

//base64加密算法
func Base64EncodeString(str string) string {
	strbytes := []byte(str)
	encoded := base64.StdEncoding.EncodeToString(strbytes)
	return encoded
}
