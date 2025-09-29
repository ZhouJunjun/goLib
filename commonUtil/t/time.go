/*
Created by
User: junjunzhou
Time: 2017/10/31 17:55
*/

package t

import (
    "strings"
    "time"
)

var (
    loc, _ = time.LoadLocation("Local")
)

// 返回年月日小时分钟秒(timestamp为0时, 默认当前时间)
func YmdHms(timestamp int64) string {
    if timestamp == 0 {
        return time.Now().Format("2006-01-02 15:04:05")
    } else {
        return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
    }
}

// 返回年月日(timestamp为0时, 默认当前时间)
func Ymd(timestamp int64) string {
    if timestamp == 0 {
        return time.Now().Format("20060102")
    } else {
        return time.Unix(timestamp, 0).Format("20060102")
    }
}

// 返回小时分钟秒(timestamp为0时, 默认当前时间)
func Hms(timestamp int64) string {
    if timestamp == 0 {
        return time.Now().Format("15:04:05")
    } else {
        return time.Unix(timestamp, 0).Format("15:04:05")
    }
}

// 第几周；
// @param sundayIsFirstDay true=周日作为每周第一天；false=周一作为每周第一天
func Week(timestamp int64, sundayIsFirstDay bool) int {
    t := time.Unix(timestamp, 0)
    
    yearDay := t.YearDay()                      // 当年第几天
    yearFirstDay := t.AddDate(0, 0, -yearDay+1) // 当年第一天
    firstDayInWeek := yearFirstDay.Weekday()    // 当年第一天是周几
    
    if sundayIsFirstDay {
        // 按周日是一周中的第一天；当年第一天是周日，当年第一天为第一周，当年第二天也是第一周
        firstWeekDays := 7 - int(firstDayInWeek) // 当年第一周有几天
        if yearDay <= firstWeekDays {
            return 1
        } else {
            if t.Weekday() == time.Saturday {
                return (yearDay-firstWeekDays)/7 + 1
            } else {
                return (yearDay-firstWeekDays)/7 + 2
            }
        }
        
    } else {
        // 按周日是一周中的第一天；当年第一天是周日，当年第一天为第一周，当年第二天是第二周
        firstWeekDays := 1 // 当年第一周有几天
        if firstDayInWeek != time.Sunday {
            firstWeekDays = 7 - int(firstDayInWeek) + 1
        }
        if yearDay <= firstWeekDays {
            return 1
        } else {
            if t.Weekday() == time.Sunday {
                return (yearDay-firstWeekDays)/7 + 1
            } else {
                return (yearDay-firstWeekDays)/7 + 2
            }
        }
    }
}

// 将格式为 yyyy-MM-dd HH:mm:ss 或 yyyyMMdd HH:mm:ss 字符串转时间戳
func Timestamp(dateTime string) int64 {
    if strings.Contains(dateTime, "-") {
        t, _ := time.ParseInLocation("2006-01-02 15:04:05", dateTime, loc)
        return t.Unix()
    } else {
        t, _ := time.ParseInLocation("20060102 15:04:05", dateTime, loc)
        return t.Unix()
    }
}

// 将格式为 yyyy-MM-dd HH:mm:ss 或 yyyyMMdd HH:mm:ss 字符串转time.Time
func Time(dateTime string) time.Time {
    if strings.Contains(dateTime, "-") {
        t, _ := time.ParseInLocation("2006-01-02 15:04:05", dateTime, loc)
        return t
    } else {
        t, _ := time.ParseInLocation("20060102 15:04:05", dateTime, loc)
        return t
    }
}

var GetCurrentTimeMillis = CurrentMs

// 获取当前毫秒时间戳
func CurrentMs() int64 {
    return time.Now().UnixNano() / 1e6
}

// 获取毫秒时间戳
func GetTimeMillis(t *time.Time) int64 {
    return t.UnixNano() / 1e6
}

// 转化mysql中读取得datetime 为 time.Time
func ParseDbDateTime(dateTime string) (time.Time, error) {
    if strings.Contains(dateTime, "+08:00") {
        dateTime = strings.Replace(dateTime, "+08:00", "", 1)
    }
    return time.ParseInLocation("2006-01-02T15:04:05", dateTime, loc)
}
