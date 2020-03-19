package whysql

import "strings"

func SafeReplaceParames(param string) string {
	var a string
	a = strings.Replace(param, "'", "‘", -1)
	//a = strings.Replace(a, "\"", "“", -1)
	a = strings.Replace(a, "&", "", -1)
	a = strings.Replace(a, "#", "", -1)
	a = strings.Replace(a, "!", "！", -1)
	a = strings.Replace(a, "*", "", -1)
	a = strings.Replace(a, "^", "", -1)
	a = strings.Replace(a, "(", "（", -1)
	a = strings.Replace(a, ")", "）", -1)
	a = strings.Replace(a, ";", "；", -1)
	a = strings.Replace(a, "--", "", -1)
	return a
}



