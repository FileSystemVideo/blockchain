package util

import (
	"fs.video/blockchain/x/copyright/config"
	"fmt"
	"github.com/shopspring/decimal"
	"math/rand"
	"strconv"
	"testing"
	"time"
)


func TestDecimal(t *testing.T) {
	dd := "5.236212212"
	ddDecimal := decimal.RequireFromString(dd)
	t.Log("number", ddDecimal.StringFixed(3))
	t.Log("number", ddDecimal.StringFixed(2))
	t.Log("number", ddDecimal.StringFixed(1))

	t.Log("number", ddDecimal.StringFixedBank(3))
	t.Log("number", ddDecimal.StringFixedBank(2))
	t.Log("number", ddDecimal.StringFixedBank(1))

	t.Log(ddDecimal.Round(2))
	t.Log(ddDecimal.RoundBank(2))

	dd = DecimalStringFixed(dd, 6)
	t.Log("data", dd)
	subDd := decimal.RequireFromString(dd)
	sjjd := ddDecimal.Sub(subDd)
	t.Log("bal", sjjd)
}

func TestRandomPublisher(t *testing.T) {
	var currentPublishId int64
	currentPublishId = 156220
	t.Log("id", BuildPublisherId(currentPublishId))
	currentPublishId = 150
	t.Log("id", BuildPublisherId(currentPublishId))
	currentPublishId = 5656210
	t.Log("id", BuildPublisherId(currentPublishId))
}


func BuildPublisherId(currentPublishId int64) string {
	if currentPublishId > config.InitPublisherId {
		currentPublishId = currentPublishId + 1
	} else {
		currentPublishId = config.InitPublisherId
	}
	return strconv.FormatInt(currentPublishId, 10)
}


var publisherIdMap map[string]string

func BuildRandomPublisherId() string {
	rand.Seed(time.Now().UnixNano())
	publisherId := rand.Int63n(900000) + 100000
	publisherIdStr := strconv.FormatInt(int64(publisherId), 10)
	if _, ok := publisherIdMap[publisherIdStr]; !ok {
		publisherIdMap[publisherIdStr] = ""
		return publisherIdStr
	}
	return ""
}

func TestRandomPublisherMap(t *testing.T) {
	publisherIdMap = make(map[string]string)
	publisherIdString := randomPublisherId(publisherIdMap)
	t.Log("id", publisherIdString)
	return

	currentTime := time.Now().Unix()
	var newMap map[string]string
	newMap = make(map[string]string)
	for i := 1000000; i < 9999999; i++ {
		publisherIdStr := strconv.Itoa(i)
		if _, ok := publisherIdMap[publisherIdStr]; ok {
			continue
		}
		newMap[publisherIdStr] = ""
		if len(newMap) > 10000 {
			endTime := time.Now().Unix()
			t.Log("time", endTime-currentTime)
			break
		}
	}
	endTime := time.Now().Unix()
	t.Log("endtime", endTime-currentTime)
}

const (
	baseRang int64 = 900000
	start    int64 = 100000
)

func randomPublisherId(publisherMap map[string]string) string {
	publisherSize := len(publisherMap)
	baseRang, start := judgePublisherStage(int64(publisherSize))
	fmt.Println("当前范围", baseRang, "************", start)
	rand.Seed(time.Now().UnixNano())
	publisherIdMap := make(map[int]string)
	i := 0
	var j int64 = 0

	for ; j < baseRang; j++ {
		publisherId := rand.Int63n(baseRang) + start
		publisherIdStr := strconv.FormatInt(publisherId, 10)
		/*if _, ok := publisherIdMap[publisherIdStr]; ok {
			continue
		}*/
		if _, ok := publisherMap[publisherIdStr]; !ok {
			publisherIdMap[i] = publisherIdStr
			i += 1
			if i > 10 {
				break
			}
		}
	}

	publisherIdString := ""
	if len(publisherIdMap) > 0 {
		publisherIndex := rand.Intn(len(publisherIdMap))
		publisherIdString = publisherIdMap[publisherIndex]
	} else {
		var k int64 = start
		for ; k < baseRang; k++ {
			publisherIdStr := strconv.FormatInt(k, 10)
			if _, ok := publisherMap[publisherIdStr]; !ok {
				publisherIdString = publisherIdStr
			}
		}
	}
	return publisherIdString
}


func judgePublisherStage(publishIdSize int64) (int64, int64) {
	var i int64 = 1
	var totalSize int64
	for true {
		if i == 1 {
			totalSize = baseRang - 1
			if totalSize > publishIdSize {
				return baseRang, start
			}
		} else {

			var square int64 = 10
			for j := 2; int64(j) < i; j++ {
				square = square * 10
			}
			currentBase := baseRang * square
			currentStart := start * square
			currentTotal := currentBase - 1
			totalSize += currentTotal
			if totalSize > publishIdSize {
				return currentBase, currentStart
			}
		}
		i += 1
	}
	return baseRang, start
}
