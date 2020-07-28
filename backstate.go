package whysql

import (
	"fmt"
)

func BackErrs(msg string) string {
	str := fmt.Sprintf("{\"state\":\"9999\",\"msg\":\"%s\"}", msg)
	return str
}

func BackErr(msg string, err error) string {
	str := fmt.Sprintf("{\"state\":\"9999\",\"msg\":\"%s\",\"err\":\"%v\"}", msg, err)
	return str
}

func BackOk(msg string) string {
	str := fmt.Sprintf("{\"state\":\"0000\",\"msg\":\"%s\"}", msg)
	return str
}

func BackOkdata(msg string, data string) string {
	str := fmt.Sprintf("{\"state\":\"0000\",\"msg\":\"%s\",\"value\":%s}", msg, data)
	return str
}

func BackOkdatajson(msg string, data string) string {
	str := fmt.Sprintf("{\"state\":\"0000\",\"msg\":\"%s\",\"value\":\"%s\"}", msg, data)
	return str
}

//返回状态集
func BackState(outResult string, total int, msg string, state string) string {
	str := fmt.Sprintf("{\"total\":%d,\"msg\":\"%s\",\"state\":\"%s\",\"rows\":%s}", total, msg, state, outResult)
	return str
}

//返回状态集带消耗时间
func BackStatetime(outResult string, total int, msg string, state string, spendtime string) string {
	str := fmt.Sprintf("{\"total\":%d,\"msg\":\"%s\",\"state\":\"%s\",\"spendtime\":\"%s\",\"rows\":%s}", total, msg, state, spendtime, outResult)
	return str
}

//返回状态 单行 结果集存放在  rows
func BackStateNo(outResult string, msg string, state string) string {
	str := fmt.Sprintf("{\"msg\":\"%s\",\"state\":\"%s\",\"rows\":%s}", msg, state, outResult)
	return str
}

//返回状态 单行 结果集存放在  data
func BackStateData(outResult string, msg string, state string) string {
	str := fmt.Sprintf("{\"msg\":\"%s\",\"state\":\"%s\",\"value\":%s}", msg, state, outResult)
	return str
}

//状态自定义返回
func BackStateInfo(msg string, state string) string {
	str := fmt.Sprintf("{\"state\":\"%s\",\"msg\":\"%s\"}", state, msg)
	return str
}

//返回状态集带消耗时间
func BackProcess(total int, rownum int, msg string, state string, data string) string {
	percentage := "0"
	if total > 0 {
		percentage = fmt.Sprintf("%.2f", float64(rownum)/float64(total)*100)
	}
	str := fmt.Sprintf("{\"total\":%d, \"rownum\":%d, \"percentage\":%s,\"msg\":\"%s\",\"state\":\"%s\",\"value\":%s}", total, rownum, percentage, msg, state, data)
	return str
}

//返回状态集带消耗时间
func BackImport(successRows int, errorRows int, msg string, state string, data string) string {
	str := fmt.Sprintf("{\"successrows\":%d, \"errorrows\":%d, \"msg\":\"%s\",\"state\":\"%s\",\"value\":%s}", successRows, errorRows, msg, state, data)
	return str
}
