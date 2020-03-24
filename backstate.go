package whysql

import (
	"fmt"
)

func BackErrs(msg string) string {
	str := fmt.Sprintf("{\"state\":\"9999\",\"msg\":\"%s\"}", msg)
	return str
}

func BackErr(msg string, err error) string {
	str := fmt.Sprintf("{\"state\":\"9999\",\"MSG\":\"%s\",\"err\":\"%v\"}", msg, err)
	return str
}

func BackOk(msg string) string {
	str := fmt.Sprintf("{\"state\":\"0000\",\"MSG\":\"%s\"}", msg)
	return str
}

func BackOkdata(msg string, data string) string {
	str := fmt.Sprintf("{\"state\":\"0000\",\"MSG\":\"%s\",\"VALUE\":%s}", msg, data)
	return str
}

func BackOkdatajson(msg string, data string) string {
	str := fmt.Sprintf("{\"state\":\"0000\",\"MSG\":\"%s\",\"VALUE\":\"%s\"}", msg, data)
	return str
}

//返回状态集
func BackState(outResult string, total int, msg string, state string) string {
	str := fmt.Sprintf("{\"TOTAL\":%d,\"MSG\":\"%s\",\"state\":\"%s\",\"ROWS\":%s}", total, msg, state, outResult)
	return str
}

//返回状态集带消耗时间
func BackStatetime(outResult string, total int, msg string, state string, spendtime string) string {
	str := fmt.Sprintf("{\"TOTAL\":%d,\"MSG\":\"%s\",\"state\":\"%s\",\"spendtime\":\"%s\",\"ROWS\":%s}", total, msg, state, spendtime, outResult)
	return str
}

//返回状态 单行 结果集存放在  rows
func BackStateNo(outResult string, msg string, state string) string {
	str := fmt.Sprintf("{\"MSG\":\"%s\",\"state\":\"%s\",\"ROWS\":%s}", msg, state, outResult)
	return str
}

//返回状态 单行 结果集存放在  data
func BackStateData(outResult string, msg string, state string) string {
	str := fmt.Sprintf("{\"MSG\":\"%s\",\"state\":\"%s\",\"VALUE\":%s}", msg, state, outResult)
	return str
}

//状态自定义返回
func BackStateInfo(msg string, state string) string {
	str := fmt.Sprintf("{\"state\":\"%s\",\"MSG\":\"%s\"}", state, msg)
	return str
}

//返回状态集带消耗时间
func BackProcess(total int, rownum int, msg string, state string, data string) string {
	percentage := "0"
	if total > 0 {
		percentage = fmt.Sprintf("%.2f", float64(rownum)/float64(total)*100)
	}
	str := fmt.Sprintf("{\"TOTAL\":%d, \"ROWNUM\":%d, \"PERCENTAGE\":%s,\"MSG\":\"%s\",\"STATE\":\"%s\",\"VALUE\":%s}", total, rownum, percentage, msg, state, data)
	return str
}

//返回状态集带消耗时间
func BackImport(successRows int, errorRows int, msg string, state string, data string) string {
	str := fmt.Sprintf("{\"SUCCESSROWS\":%d, \"ERRORROWS\":%d, \"MSG\":\"%s\",\"state\":\"%s\",\"VALUE\":%s}", successRows, errorRows, msg, state, data)
	return str
}
