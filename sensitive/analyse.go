package sensitive

import (
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)



var SenSitiveStatus int = 1 
//var PersonSensitiveMap = make(map[string]interface{}) 
//var systemSensitiveMap = make(map[string]interface{}) 
var SensitiveWordSplit = "|*|"
var SensitiveStatusSplit = "|-|"
var SysSenStatus = 1 
var PerSenStatus = 1 

/*
	
    
    
	()
*/
func SensitiveAnalyse(word string) (bool, int, int) {
	if word == "" {
		return false, 0, 0
	}
	
	if systemSensitive == nil {
		return false, 0, 0
	}
	_, flag := systemSensitive.CheckSensitive(word)
	if flag {
		return true, 1, 0
	}
	/*for k, _ := range systemSensitiveMap {
		if strings.Contains(word, k) {

		}
	}*/
	
	if personSensitive != nil {
		_, flag = personSensitive.CheckSensitive(word)
		if flag {
			return true, 0, 1
		}
	}
	return false, 0, 0
}


func SensitiveWholeAnalyse(word string) (int, int) {
	if word == "" {
		return 0, 0
	}
	systemStatus := 0
	personStatus := 0
	
	if systemSensitive == nil {
		return systemStatus, PerSenStatus
	}
	_, flag := systemSensitive.CheckSensitive(word)
	if flag {
		systemStatus = 1
	}
	
	if personSensitive != nil {
		_, flag = personSensitive.CheckSensitive(word)
		if flag {
			personStatus = 1

		}
	}
	return systemStatus, personStatus
}

/*
	
	,,
*/
func PersonSensitiveAnalyse(word string) int {
	if word == "" {
		return 0
	}
	
	_, flag := personSensitive.CheckSensitive(word)
	if flag {
		return 1
	}
	//for i := 0; i < len(PersonSensitiveMap); i++ {
	/*for k, _ := range PersonSensitiveMap {
		if strings.Contains(word, k) {
			return 1
		}
	}*/
	return 0
}


func SystemSensitiveAnalyse(word string) int {
	if word == "" {
		return 0
	}
	
	_, flag := systemSensitive.CheckSensitive(word)
	_, flag = personSensitive.CheckSensitive(word)
	if flag {
		return 1
	}

	return 0
}


func personSensitiveWords() {
	personWords, err := QuerySensitive()
	if err != nil {
		logrus.Error("", err)
		//return
	}
	if personWords != "" {
		sensitiveArray := strings.Split(personWords, SensitiveStatusSplit)
		statusString := sensitiveArray[0]
		status, err := strconv.Atoi(statusString)
		if err != nil {
			logrus.Error("", err)
		} else {
			SenSitiveStatus = status
		}
		if len(sensitiveArray) > 1 && sensitiveArray[1] != "" {
			personWords := strings.Replace(sensitiveArray[1], SensitiveWordSplit, " ", -1)
			initPersonSensitiveWords(strings.Fields(personWords))
			/*wordsArray := strings.Split(sensitiveArray[1], SensitiveWordSplit)
			for i := 0; i < len(wordsArray); i++ {
				PersonSensitiveMap[wordsArray[i]] = nil
			}*/
		}
	}
}


func PersonSensitiveWordsUpdate(newpersonword string) {
	if newpersonword != "" {
		initPersonSensitiveWords(strings.Fields(newpersonword))
	} else {
		personSensitive = initSensitiveMap()
	}
}


/*func initSystemSensitiveWords() {
	wordsArray := strings.Split(systemword, ",")
	for i := 0; i < len(wordsArray); i++ {
		systemSensitiveMap[wordsArray[i]] = nil
	}
}*/
