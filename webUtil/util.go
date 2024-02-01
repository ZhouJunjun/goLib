/*
Created by
User: junjunzhou
Time: 2017/9/7 15:02
*/

package webUtil

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// 获取客户端ip
func GetIp(request *http.Request) string {

	if ip := strings.TrimSpace(request.Header.Get("X-Forwarded-For")); ip != "" && ip != "unKnown" {
		ips := strings.Split(ip, ",")
		if len(ips) > 0 && ips[0] != "" {
			return ips[0]
		}
	}

	if ip := strings.TrimSpace(request.Header.Get("X-Real-IP")); ip != "" && ip != "unKnown" {
		return ip
	}

	return strings.Split(request.RemoteAddr, ":")[0]
}

// 获取cookie
func GetCookie(request *http.Request, name string) string {
	cookie, err := request.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}

// 读取http表单参数, 返回string类型
func GetParamString(request *http.Request, key string) string {
	arr, ok := request.Form[key]
	if !ok || len(arr) <= 0 {
		return ""
	}
	return strings.TrimSpace(arr[0])
}

// 读取http表单参数, 返回string类型
func GetParamStringD(request *http.Request, key string, defaultVal string) string {
	arr, ok := request.Form[key]
	if !ok || len(arr) <= 0 {
		return defaultVal
	}
	return strings.TrimSpace(arr[0])
}

// 读取http表单参数, 返回int64类型
func GetParamInt64(request *http.Request, key string) int64 {
	arr, ok := request.Form[key]
	if !ok || len(arr) <= 0 {
		return 0
	}
	n, _ := strconv.ParseInt(arr[0], 10, 64)
	return n
}

// 读取http表单参数, 返回int64类型
func GetParamInt64D(request *http.Request, key string, defaultVal int64) int64 {
	arr, ok := request.Form[key]
	if !ok || len(arr) <= 0 || len(arr[0]) == 0 {
		return defaultVal
	}
	n, _ := strconv.ParseInt(arr[0], 10, 64)
	return n
}

// 读取http表单参数, 返回int类型
func GetParamInt(request *http.Request, key string) int {
	arr, ok := request.Form[key]
	if !ok || len(arr) <= 0 {
		return 0
	}
	n, _ := strconv.Atoi(arr[0])
	return n
}

// 读取http表单参数, 返回int类型
func GetParamIntD(request *http.Request, key string, defaultVal int) int {
	arr, ok := request.Form[key]
	if !ok || len(arr) <= 0 || len(arr[0]) == 0 {
		return defaultVal
	}
	n, _ := strconv.Atoi(arr[0])
	return n
}

// 表单参数按字段名排序转字符串
func SortQuery(request *http.Request, passKey string) string {

	keys := []string{}
	for key := range request.Form {
		if key != passKey {
			keys = append(keys, key)
		}
	}

	quickSort(keys)

	buf := bytes.Buffer{}
	for _, key := range keys {
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(key)
		buf.WriteByte('=')
		buf.WriteString(request.Form.Get(key))
	}
	return buf.String()
}

func quickSort(keys []string) {
	if len(keys) <= 1 {
		return
	}

	mid, i := keys[0], 1
	head, tail := 0, len(keys)-1

	for head < tail {
		if keys[i] > mid {
			keys[i], keys[tail] = keys[tail], keys[i]
			tail--
		} else {
			keys[i], keys[head] = keys[head], keys[i]
			head++
			i++
		}
	}
	keys[head] = mid

	quickSort(keys[:head])
	quickSort(keys[head+1:])
}

var (
	ipv4 = regexp.MustCompile(`^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$`)
)

// 获取hostname对应的ip
func GetHostNameIp() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("os.Hostname error, %s", err.Error())
	} else if hostname == "" {
		return "", fmt.Errorf("hostname is ''")
	}

	addrSlice, err := net.LookupHost(hostname)
	if err != nil {
		return "", fmt.Errorf("net.LookupHost error, %s", err.Error())
	}
	if len(addrSlice) <= 0 {
		return "", fmt.Errorf("hostname ip not found")
	}

	// 优先输出ipv4格式ip
	for _, addr := range addrSlice {
		if ipv4.MatchString(addr) {
			return addr, nil
		}
	}
	return addrSlice[0], nil
}
