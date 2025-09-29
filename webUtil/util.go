/*
Created by
User: junjunzhou
Time: 2017/9/7 15:02
*/

package webUtil

import (
    "bytes"
    "encoding/json"
    "fmt"
    "github.com/ZhouJunjun/goLib/commonUtil/s"
    "github.com/ZhouJunjun/goLib/commonUtil/t"
    "github.com/ZhouJunjun/goLib/util"
    "github.com/ZhouJunjun/goLib/commonUtil/logBuffer"
    "net"
    "net/http"
    "net/url"
    "os"
    "regexp"
    "sort"
    "strconv"
    "strings"
)

// 获取请求客户端ip
func GetIp(request *http.Request) string {

    if ip := strings.TrimSpace(request.Header.Get("X-Forwarded-For")); ip != "" && ip != "unKnown" {
        ips := strings.Split(ip, ",")
        if len(ips) > 0 && ips[0] != "" {
            return strings.TrimSpace(ips[0])
        }
    }

    if ip := strings.TrimSpace(request.Header.Get("X-Real-IP")); ip != "" && ip != "unKnown" {
        return ip
    }

    ip := strings.Split(request.RemoteAddr, ":")[0]
    return strings.TrimSpace(ip)
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
    return GetParamStringD(request, key, "")
}

// 读取http表单参数, 返回string类型
func GetParamStringD(request *http.Request, key string, defaultVal string) string {
    arr, ok := request.Form[key]
    if !ok || len(arr) <= 0 {
        return defaultVal
    }
    return strings.TrimSpace(arr[0])
}

// 读取http表单参数, 返回*string类型, 返回nil=未传值
func GetParamStringN(request *http.Request, key string) *string {
    arr, ok := request.Form[key]
    if !ok || len(arr) <= 0 {
        return nil
    }
    s := strings.TrimSpace(arr[0])
    return &s
}

// 读取http表单参数, 返回int64类型
func GetParamInt64(request *http.Request, key string) int64 {
    return GetParamInt64D(request, key, 0)
}

// 读取http表单参数, 返回int64类型
func GetParamInt64D(request *http.Request, key string, defaultVal int64) int64 {
    arr, ok := request.Form[key]
    if !ok || len(arr) <= 0 {
        return defaultVal
    }

    if arr[0] = strings.TrimSpace(arr[0]); s.IsNumber(arr[0]) {
        n, _ := strconv.ParseInt(arr[0], 10, 64)
        return n
    }
    return defaultVal
}

// 读取http表单参数, 返回int64类型, 返回nil=未传数值
func GetParamInt64N(request *http.Request, key string) *int64 {
    arr, ok := request.Form[key]
    if !ok || len(arr) <= 0 {
        return nil
    }

    if arr[0] = strings.TrimSpace(arr[0]); s.IsNumber(arr[0]) {
        n, _ := strconv.ParseInt(arr[0], 10, 64)
        return &n
    }
    return nil
}

// 读取http表单参数, 返回int类型
func GetParamInt(request *http.Request, key string) int {
    return GetParamIntD(request, key, 0)
}

// 读取http表单参数, 返回int类型
func GetParamIntD(request *http.Request, key string, defaultVal int) int {
    arr, ok := request.Form[key]
    if !ok || len(arr) <= 0 {
        return defaultVal
    }

    if arr[0] = strings.TrimSpace(arr[0]); s.IsNumber(arr[0]) {
        n, _ := strconv.Atoi(arr[0])
        return n
    }
    return defaultVal
}

// 读取http表单参数, 返回int类型
func GetParamIntN(request *http.Request, key string) *int {
    arr, ok := request.Form[key]
    if !ok || len(arr) <= 0 {
        return nil
    }

    if arr[0] = strings.TrimSpace(arr[0]); s.IsNumber(arr[0]) {
        n, _ := strconv.Atoi(arr[0])
        return &n
    }
    return nil
}

// 读取http表单参数, 返回float32类型
func GetParamFloat32(request *http.Request, key string) float32 {
    return GetParamFloat32D(request, key, 0)
}

// 读取http表单参数, 返回float32类型
func GetParamFloat32D(request *http.Request, key string, defaultVal float32) float32 {
    arr, ok := request.Form[key]
    if !ok || len(arr) <= 0 {
        return defaultVal
    }

    if arr[0] = strings.TrimSpace(arr[0]); s.IsNumber(arr[0]) {
        n, _ := strconv.ParseFloat(arr[0], 32)
        return float32(n)
    }
    return defaultVal
}

// 读取http表单参数, 返回float64类型
func GetParamFloat64(request *http.Request, key string) float64 {
    return GetParamFloat64D(request, key, 0)
}

// 读取http表单参数, 返回float64类型
func GetParamFloat64D(request *http.Request, key string, defaultVal float64) float64 {
    arr, ok := request.Form[key]
    if !ok || len(arr) <= 0 {
        return defaultVal
    }

    if arr[0] = strings.TrimSpace(arr[0]); s.IsNumber(arr[0]) {
        n, _ := strconv.ParseFloat(arr[0], 64)
        return n
    }
    return defaultVal
}

// 读取http表单参数, 返回string类型
func GetPostString(request *http.Request, key string) string {
    return GetPostStringD(request, key, "")
}

// 读取http表单参数, 返回string类型
func GetPostStringD(request *http.Request, key string, defaultVal string) string {
    arr, ok := request.PostForm[key]
    if !ok || len(arr) <= 0 {
        return defaultVal
    }
    return strings.TrimSpace(arr[0])
}

// 读取http表单参数, 返回int64类型
func GetPostInt64(request *http.Request, key string) int64 {
    return GetPostInt64D(request, key, 0)
}

// 读取http表单参数, 返回int64类型
func GetPostInt64D(request *http.Request, key string, defaultVal int64) int64 {
    arr, ok := request.PostForm[key]
    if !ok || len(arr) <= 0 {
        return defaultVal
    }

    if arr[0] = strings.TrimSpace(arr[0]); s.IsNumber(arr[0]) {
        n, _ := strconv.ParseInt(arr[0], 10, 64)
        return n
    }
    return defaultVal
}

// 读取http表单参数, 返回int类型
func GetPostInt(request *http.Request, key string) int {
    return GetPostIntD(request, key, 0)
}

// 读取http表单参数, 返回int类型
func GetPostIntD(request *http.Request, key string, defaultVal int) int {
    arr, ok := request.PostForm[key]
    if !ok || len(arr) <= 0 {
        return defaultVal
    }

    if arr[0] = strings.TrimSpace(arr[0]); s.IsNumber(arr[0]) {
        n, _ := strconv.Atoi(arr[0])
        return n
    }
    return defaultVal
}

func GetParamJsonTime(request *http.Request, key string) *t.JsonTime {
    return GetParamJsonTimeD(request, key, nil)
}

func GetParamJsonTimeD(request *http.Request, key string, defaultVal *t.JsonTime) *t.JsonTime {
    arr, ok := request.Form[key]
    if !ok || len(arr) <= 0 {
        return defaultVal
    }

    t := &t.JsonTime{}
    if err := json.Unmarshal([]byte(arr[0]), t); err == nil {
        return t
    }
    return defaultVal
}

func GetPostJsonTime(request *http.Request, key string) *t.JsonTime {
    return GetPostJsonTimeD(request, key, nil)
}

func GetPostJsonTimeD(request *http.Request, key string, defaultVal *t.JsonTime) *t.JsonTime {
    arr, ok := request.PostForm[key]
    if !ok || len(arr) <= 0 {
        return defaultVal
    }

    t := &t.JsonTime{}
    if err := json.Unmarshal([]byte(arr[0]), t); err == nil {
        return t
    }
    return defaultVal
}

// 表单参数按字段名排序转字符串: a=A&b=B
// @param urlValues 取自request.Form or request.PostForm
// @param excludeKeys 需要排除的字段名
func GetSortQuery(urlValues url.Values, excludeKeys ...string) string {

    keySlice, i := make([]string, len(urlValues)), 0
    for key := range urlValues {
        if excludeKeys == nil || !util.InSliceString(excludeKeys, key) {
            keySlice[i], i = key, i+1
        }
    }

    keySlice = keySlice[:i]
    quickSort(keySlice)

    buf := bytes.Buffer{}
    for _, key := range keySlice {
        if buf.Len() > 0 {
            buf.WriteByte('&')
        }
        buf.WriteString(key)
        buf.WriteByte('=')
        buf.WriteString(urlValues.Get(key))
    }
    return buf.String()
}

// 表单参数按字段名排序转字符串
// Deprecated: use GetSortQuery(...) instead
func SortQuery(request *http.Request, excludeKey string) string {
    return GetSortQuery(request.Form, excludeKey)
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

// 获取服务端本地ip
func GetLocalServerIp() (string, error) {
    if iFaceSlice, err := net.Interfaces(); err == nil {
        if iFaceSlice != nil {
            for _, iFace := range iFaceSlice {
                if iFace.Flags&net.FlagLoopback != 0 {
                    continue // loopback interface
                }
                if iFace.Flags&net.FlagUp == 0 {
                    continue // interface down
                }

                if tmpAddrSlice, err := iFace.Addrs(); err == nil && tmpAddrSlice != nil {
                    for _, address := range tmpAddrSlice {
                        if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
                            if ipNet.IP.To4() != nil {
                                return ipNet.IP.String(), nil
                            }
                        }
                    }
                }
            }
        }
        return "", fmt.Errorf("ip not found")
    } else {
        return "", err
    }

    // 不使用 net.Interfaces()，不能过滤带Loopback属性的网络设备绑定的ip；ipNet.IP.IsLoopback()中的仅过滤了 127.* 和 ::1
    /*interfaceSlice, err := net.Interfaces()
      if err != nil {
          return "", err
      }

      for _, each := range interfaceSlice {
          if each.Name == "eth0" {
              if addrSlice, err := each.Addrs(); err != nil {
                  return "", err
              } else {
                  for _, address := range addrSlice {
                      if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
                          if ipNet.IP.To4() != nil {
                              return ipNet.IP.String(), nil
                          }
                      }
                  }
                  return "", fmt.Errorf("%s's ip not found", each.Name)
              }
          }
      }
      return "", fmt.Errorf("eth0 not found")*/
}

var ipv4 = regexp.MustCompile(`^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$`)

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

// 标准格式request log；包含url、postBody、登录userId
func GetRequestLog(request *http.Request) *logBuffer.Buffer {

    buffer := new(logBuffer.Buffer).Append("url:").Append(request.URL.Path)

    if len(request.URL.RawQuery) > 0 {
        buffer.Append("?").Append(request.URL.RawQuery)
    }

    if l := len(request.PostForm); l > 0 {
        buffer.Append(", body:{")

        keySlice, i := make([]string, l), 0
        for key := range request.PostForm {
            keySlice[i], i = key, i+1
        }
        sort.Slice(keySlice, func(i, j int) bool {
            return keySlice[i] < keySlice[j]
        })

        for _, key := range keySlice {
            slice := request.PostForm[key]
            if len(slice) == 1 {
                buffer.AppendFormat("%s:%s, ", key, slice[0])
            } else {
                buffer.AppendFormat("%s:%v, ", key, slice)
            }
        }
        buffer.Del(-2, 2).Append("}")
    }
    if cookie, _ := request.Cookie("member_id"); cookie != nil {
        buffer.Append(", userId:").Append(strings.TrimSuffix(cookie.Value, "%4056.com"))
    }
    return buffer
}

func GetHeader(request *http.Request, name string) *string {
    if request != nil && request.Header != nil {
        if values, ok := request.Header[name]; ok {
            return &values[0]
        }
    }
    return nil
}
