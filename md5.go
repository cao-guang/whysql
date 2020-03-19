package whysql

import (
	"crypto/md5"
	"strings"
	"encoding/hex"
	"fmt"
)

//条件值加密成redis key
func GetKeyWithParam(ip_key_hctj string) string {
	var vsKey string
	vsKey = fmt.Sprintf("%s", ip_key_hctj)
	return MD5(vsKey, true)
}

func MD5(msg string, upper bool) string {
	h := md5.New()
	h.Write([]byte(msg))
	cipherStr := h.Sum(nil)
	if upper {
		return strings.ToUpper(hex.EncodeToString(cipherStr))
	}
	return hex.EncodeToString(cipherStr)
}


