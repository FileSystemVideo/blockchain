package sensitive

import (
	"io/ioutil"
	"os"
	"strings"
)

type SensitiveMap struct {
	sensitiveNode map[string]interface{}
	isEnd         bool
}

var personSensitive *SensitiveMap
var systemSensitive *SensitiveMap


func initSystemSensitiveWords() *SensitiveMap {
	if systemSensitive == nil {
		systemSensitive = InitDictionary(systemSensitive)
	}
	return systemSensitive
}


func initPersonSensitiveWords(personWords []string) *SensitiveMap {
	personSensitive = InitPersonDictionary(personSensitive, personWords)
	return personSensitive
}


func initSensitiveMap() *SensitiveMap {
	return &SensitiveMap{
		sensitiveNode: make(map[string]interface{}),
		isEnd:         false,
	}
}


func ClearPersonSensitive() {
	personSensitive = &SensitiveMap{
		sensitiveNode: make(map[string]interface{}),
		isEnd:         false,
	}
}


func readDictionary(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	str, err := ioutil.ReadAll(file)
	dictionary := strings.Fields(string(str))
	return dictionary
}


func InitDictionary(s *SensitiveMap) *SensitiveMap {
	s = initSensitiveMap()

	dictionary := strings.Fields(systemword)
	for _, words := range dictionary {
		sMapTmp := s
		w := []rune(words)
		wordsLength := len(w)
		for i := 0; i < wordsLength; i++ {
			t := string(w[i])
			isEnd := false

			if i == (wordsLength - 1) {
				isEnd = true
			}
			func(tx string) {
				if _, ok := sMapTmp.sensitiveNode[tx]; !ok {
					sMapTemp := new(SensitiveMap)
					sMapTemp.sensitiveNode = make(map[string]interface{})
					sMapTemp.isEnd = isEnd
					sMapTmp.sensitiveNode[tx] = sMapTemp
				}
				sMapTmp = sMapTmp.sensitiveNode[tx].(*SensitiveMap)
				sMapTmp.isEnd = isEnd
			}(t)
		}
	}
	return s
}


func InitPersonDictionary(s *SensitiveMap, personWords []string) *SensitiveMap {
	s = initSensitiveMap()
	for _, words := range personWords {
		sMapTmp := s
		w := []rune(words)
		wordsLength := len(w)
		for i := 0; i < wordsLength; i++ {
			t := string(w[i])
			isEnd := false

			if i == (wordsLength - 1) {
				isEnd = true
			}
			func(tx string) {
				if _, ok := sMapTmp.sensitiveNode[tx]; !ok {
					sMapTemp := new(SensitiveMap)
					sMapTemp.sensitiveNode = make(map[string]interface{})
					sMapTemp.isEnd = isEnd
					sMapTmp.sensitiveNode[tx] = sMapTemp
				}
				sMapTmp = sMapTmp.sensitiveNode[tx].(*SensitiveMap)
				sMapTmp.isEnd = isEnd
			}(t)
		}
	}
	return s
}


func (s *SensitiveMap) CheckSensitive(text string) (string, bool) {
	text = strings.ReplaceAll(text, " ", "")
	content := []rune(text)
	contentLength := len(content)
	result := false
	ta := ""
	for index := range content {
		sMapTmp := s
		target := ""
		in := index
		for {
			wo := string(content[in])
			target += wo
			if _, ok := sMapTmp.sensitiveNode[wo]; ok {
				if sMapTmp.sensitiveNode[wo].(*SensitiveMap).isEnd {
					result = true
					break
				}
				if in == contentLength-1 {
					break
				}
				sMapTmp = sMapTmp.sensitiveNode[wo].(*SensitiveMap)
				in++
			} else {
				break
			}
		}
		if result {
			ta = target
			break
		}
	}
	return ta, result
}


type Target struct {
	Indexes []int
	Len     int
}

func (s *SensitiveMap) FindAllSensitive(text string) map[string]*Target {
	content := []rune(text)
	contentLength := len(content)
	result := false

	ta := make(map[string]*Target)
	for index := range content {
		sMapTmp := s
		target := ""
		in := index
		result = false
		for {
			wo := string(content[in])
			target += wo
			if _, ok := sMapTmp.sensitiveNode[wo]; ok {
				if sMapTmp.sensitiveNode[wo].(*SensitiveMap).isEnd {
					result = true
					break
				}
				if in == contentLength-1 {
					break
				}
				sMapTmp = sMapTmp.sensitiveNode[wo].(*SensitiveMap)
				in++
			} else {
				break
			}
		}
		if result {
			if _, targetInTa := ta[target]; targetInTa {
				ta[target].Indexes = append(ta[target].Indexes, index)
			} else {
				ta[target] = &Target{
					Indexes: []int{index},
					Len:     len([]rune(target)),
				}
			}
		}
	}
	return ta
}
