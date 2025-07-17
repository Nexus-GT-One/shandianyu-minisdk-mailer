package timeutil

import (
	"fmt"
	"time"
)

// 同时输出指定时间的中国时间和美国时间
//
// specifiedTime-指定时间
func ShowChinaTimeAndAmericaTime(specifiedTime time.Time) string {
	// 获取当前北京时间
	beijingLocation, _ := time.LoadLocation("Asia/Shanghai")
	beijingTime := specifiedTime.In(beijingLocation)

	// 转换为洛杉矶时间
	losAngelesLocation, _ := time.LoadLocation("America/Los_Angeles")
	losAngelesTime := specifiedTime.In(losAngelesLocation)

	return fmt.Sprintf("【中国时间：%s (%s) ~ 美国时间：%s (%s)】", beijingTime.Format(time.DateTime), beijingTime.Weekday(), losAngelesTime.Format(time.DateTime), losAngelesTime.Weekday())
}

// 获取指定时间所在当天的 UTC 0点时间戳
func DateStart(t time.Time) time.Time {
	dateStart := fmt.Sprintf("%s 00:00:00", time.Unix(t.Unix(), 0).Format(time.DateOnly))
	_time, _ := time.Parse(time.DateTime, dateStart)
	return time.Date(_time.Year(), _time.Month(), _time.Day(), 0, 0, 0, 0, time.UTC)
}

// 获取当前年份1月1日凌晨0点的时间戳
func GetYearStart() int64 {
	// 获取当前时间
	currentTime := time.Now()

	// 获取当前年份
	year := currentTime.Year()

	// 创建当前年份1月1日的0点0分0秒时间
	startOfYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local)

	// 获取该时间点的时间戳
	return startOfYear.UnixMilli()
}
