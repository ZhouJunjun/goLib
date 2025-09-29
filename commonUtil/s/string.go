/**
 * @author junjunzhou
 * @date 2024/10/29
 */
package s

import (
    "bytes"
    "fmt"
    "github.com/ZhouJunjun/goLib/log4j"
    "regexp"
    "strconv"
)

var (
    // 判断是否是纯数字
    numReg, _  = regexp.Compile(`^-?[0-9]+(\.[0-9]+)?$`)
    pIntReg, _ = regexp.Compile(`^[0-9]+$`)
)

// 判断是否是数字, 含负数,正数,小数
func IsNumber(s string) bool {
    return numReg.MatchString(s)
}

// 判断是否是正整数(positive integer)
func IsPInt(s string) bool {
    return pIntReg.MatchString(s)
}

// 判断是否为正整数拼接; delimiter为分割符
func IsPIntJoin(s string, delimiter string) bool {

    size, dSize := len(s), len(delimiter)
    if size == 0 || dSize == 0 {
        return false
    }

    for i := 0; i < size; i++ {
        if c := s[i]; c >= '0' && c <= '9' {
            // ok
        } else if i == 0 { // 非数字开头
            return false

        } else { // 判断分割符
            if i+dSize >= size {
                return false
            }

            for j := 0; j < dSize; j++ { // 比较分割符
                if s[i+j] != delimiter[j] {
                    return false
                }
            }

            i += dSize - 1 // 分割符相等
        }
    }
    return true
}

// 字符串转int64
func ToInt64(s string) int64 {
    i, err := strconv.ParseInt(s, 10, 64)
    if err != nil {
        _ = log4j.ErrorStack("'%s' to int 64 error: %s", s, err.Error())
        return 0
    }
    return i
}

// 字符串转int
func ToInt(s string) int {
    i, err := strconv.Atoi(s)
    if err != nil {
        _ = log4j.ErrorStack("'%s' to int error: %s", s, err.Error())
        return 0
    }
    return i
}

// 字符串转布尔类型
func ToBool(s string) bool {
    b, err := strconv.ParseBool(s)
    if err != nil {
        _ = log4j.ErrorStack("'%s' to bool error:", s, err.Error())
        return false
    }
    return b
}

// 转小写
func ToLower(s string) string {
    if len(s) <= 0 {
        return s
    }

    // 65-90 -> 97-122
    bf, bs := bytes.Buffer{}, []byte(s)
    for _, c := range bs {
        if c >= 'A' && c <= 'Z' {
            bf.WriteByte(c + 32)
        } else {
            bf.WriteByte(c)
        }
    }
    return bf.String()
}

// int64 转 string
func ValOfInt64(num int64) string {
    return fmt.Sprintf("%d", num)
}

// int 转 string
func ValOfInt(num int) string {
    return fmt.Sprintf("%d", num)
}

// boot 转 string
func ValOfBool(f bool) string {
    if f {
        return "true"
    }
    return "false"
}

// 反转字符串
func Reverse(target string) string {
    runes := []rune(target)
    for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
        runes[from], runes[to] = runes[to], runes[from]
    }
    return string(runes)
}
