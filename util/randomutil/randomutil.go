package randomutil

import (
	"math/rand"
)

// 生成随机字符串
//
// length-字符串长度
func String(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	for i := 0; i < length; i++ {
		result = append(result, bytes[Int(0, len(bytes)-1)])
	}
	return string(result)
}

// 生成随机整数
//
// min-最小值
// max-最大值
func Int(min, max int) int {
	return rand.Intn(max-min+1) + min
}

// 生成64位随机整数
//
// min-最小值
// max-最大值
func Int64(min, max int64) int64 {
	return rand.Int63n(max-min+1) + min
}
