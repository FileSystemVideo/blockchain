package util

import (
	"fs.video/blockchain/core"
	"regexp"
	"strings"
)

/*
	decimal ,
*/
func DecimalStringFixed(originNum string, precision int) string {
	if strings.Contains(originNum, ".") {
		originArray := strings.Split(originNum, ".")
		if len(originArray[1]) > precision {
			price := string([]rune(originArray[1])[:precision])
			originNum = originArray[0] + "." + price
		}
	}
	return originNum
}

var reg1 = regexp.MustCompile(`^\d+(\.\d{1,18})?$`)

/**

*/
func JudgeAmount(realAmount string) bool {
	log := core.BuildLog(core.GetFuncName(), core.LmChainUtil)
	flag := reg1.MatchString(realAmount)
	if !flag {
		return flag
	}
	
	realArray := strings.Split(realAmount, ".")
	realIntString := realArray[0]
	if len(realArray) > 1 {
		realIntString = realArray[1]
	}
	if len(realIntString) > 18 { 
		log.WithField("amount", realAmount).Error("Amount exceeds maximum limit")
		return false
	}
	return true
}
