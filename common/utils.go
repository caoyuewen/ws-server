package common

import (
	"math/rand"
	"time"
)

var count int64
// 获取一个随机数
func RandNum(startNum, endNum int) int {
	count++
	if count >= 1<<8 {
		count = 0
	}
	rand.Seed(time.Now().UnixNano() + count)
	rnd := rand.Intn(endNum - startNum)
	return rnd + startNum
}
