/**
 * @author junjunzhou
 * @date 2021/6/10
 */
package util

import (
    "code.sohuno.com/qianfan-go/qf-library/commonUtil/s"
    "code.sohuno.com/qianfan-go/qf-library/log4j"
    "encoding/json"
    "errors"
    "fmt"
    "reflect"
    "strings"
    "time"
    "unsafe"
)

var (
    typeOfTime = reflect.TypeOf(time.Time{})
)

func If(flag bool, trueVal, falseVal interface{}) interface{} {
    if flag {
        return trueVal
    }
    return falseVal
}

func IfInt(flag bool, trueVal, falseVal int) int {
    if flag {
        return trueVal
    }
    return falseVal
}

func IfInt64(flag bool, trueVal, falseVal int64) int64 {
    if flag {
        return trueVal
    }
    return falseVal
}

func IfString(flag bool, trueVal, falseVal string) string {
    if flag {
        return trueVal
    }
    return falseVal
}

func JsonEncode(obj interface{}) string {
    bs, err := json.Marshal(obj)
    if err != nil {
        _ = log4j.ErrorStack("json marshal error, data=%+v, err=%+v", obj, err)
    }
    return string(bs)
}

func JsonDecode(str string, obj interface{}) {
    err := json.Unmarshal([]byte(str), obj)
    if err != nil {
        _ = log4j.ErrorStack("json unmarshal error, data=%+v, err=%+v", obj, err)
    }
}

func JoinInt(slice []int, sep string) string {
    
    size := len(slice)
    if size == 0 {
        return ""
    } else if size == 1 {
        return s.ValOfInt(slice[0])
    }
    
    builder := strings.Builder{}
    builder.WriteString(s.ValOfInt(slice[0]))
    
    for i := 1; i < size; i++ {
        builder.WriteString(sep)
        builder.WriteString(s.ValOfInt(slice[i]))
    }
    return builder.String()
}

func JoinInt64(slice []int64, sep string) string {
    
    size := len(slice)
    if size == 0 {
        return ""
    } else if size == 1 {
        return s.ValOfInt64(slice[0])
    }
    
    builder := strings.Builder{}
    builder.WriteString(s.ValOfInt64(slice[0]))
    
    for i := 1; i < size; i++ {
        builder.WriteString(sep)
        builder.WriteString(s.ValOfInt64(slice[i]))
    }
    return builder.String()
}

func SplitInt64(str, sep string) []int64 {
    slice := strings.Split(str, sep)
    rsSlice := make([]int64, len(slice))
    for i, each := range slice {
        rsSlice[i] = s.ToInt64(each)
    }
    return rsSlice
}

func SplitInt(str, sep string) []int {
    slice := strings.Split(str, sep)
    rsSlice := make([]int, len(slice))
    for i, each := range slice {
        rsSlice[i] = s.ToInt(each)
    }
    return rsSlice
}

func InSliceInt(slice []int, i int) bool {
    if slice == nil {
        return false
    }
    
    for _, each := range slice {
        if each == i {
            return true
        }
    }
    return false
}

func InSliceInt64(slice []int64, i int64) bool {
    if slice == nil {
        return false
    }
    
    for _, each := range slice {
        if each == i {
            return true
        }
    }
    return false
}

func InSliceString(slice []string, i string) bool {
    if slice == nil {
        return false
    }
    
    for _, each := range slice {
        if each == i {
            return true
        }
    }
    return false
}

/*
 * 判断元素是否在slice中, 支持元素类型为: 基础数据类型, 指针, slice, map
 *  user1 := &User{Name: "aaa"}
 *  user2 := &User{Name: "aaa"}
 *  userSlice1 := []*User{user1}
 *  userSlice2 := []interface{}{user1}
 *  fmt.Println(InSlice(userSlice1, user1)) //返回true
 *  fmt.Println(InSlice(userSlice1, user2)) //返回false
 *  fmt.Println(InSlice(userSlice2, user1)) //返回true
 *  fmt.Println(InSlice(userSlice2, user2)) //返回false
 */
func InSlice(slice interface{}, item interface{}) bool {
    sliceTyp, sliceVal := reflect.TypeOf(slice), reflect.ValueOf(slice)
    if sliceTyp.Kind() != reflect.Slice {
        return false
    }
    
    itemTyp, itemVal := reflect.TypeOf(item), reflect.ValueOf(item)
    
    for i, size := 0, sliceVal.Len(); i < size; i++ {
        sliceItem := sliceVal.Index(i)
        
        // slice为[]interface{}, 获取元素具体类型
        if sliceItem.Kind() == reflect.Interface {
            sliceItem = sliceItem.Elem()
        }
        
        if sliceItem.Kind() != itemTyp.Kind() {
            continue
        }
        
        switch sliceItem.Kind() {
        // 基础数据类型, 值相等(不是同一块内存) 就相当于存在
        case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
            if sliceItem.Int() == itemVal.Int() {
                return true
            }
        case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
            if sliceItem.Uint() == itemVal.Uint() {
                return true
            }
        case reflect.Float32, reflect.Float64:
            if sliceItem.Float() == itemVal.Float() {
                return true
            }
        case reflect.Bool:
            if sliceItem.Bool() == itemVal.Bool() {
                return true
            }
        
        case reflect.Ptr, reflect.Slice, reflect.Map:
            if sliceItem.Pointer() == itemVal.Pointer() {
                return true
            }
        
        case reflect.String:
            if sliceItem.String() == itemVal.String() {
                return true
            }
            
            // 如果是struct slice数值, 每一次赋值都会重新分配内存并复制整个结构体, 所以100%不存在;
            // 如果要求struct相同, 需要深度比较, 太耗费性能, 先不搞
        case reflect.Struct:
            panic("unsupported item type: " + itemTyp.Kind().String())
        
        default:
            panic("unsupported item type: " + itemTyp.Kind().String())
        }
    }
    return false
}

/* 返回map的key slice, slice的数据类型由key的类型决定; map[string]xxx返回[]string, map[int]xxx返回[]int
 * 示例：targetMap := map[string]int{}, keySlice := MapKeySlice(targetMap).([]string)
 */
func MapKeySlice(m interface{}) interface{} {
    typ, val := reflect.TypeOf(m), reflect.ValueOf(m)
    if typ.Kind() != reflect.Map {
        return nil
    }
    
    keys := val.MapKeys()
    slice := reflect.New(reflect.SliceOf(typ.Key())).Elem()
    for _, key := range keys {
        slice.Set(reflect.Append(slice, key))
    }
    return slice.Interface()
}

/*
 * 返回map的val slice, slice的数据类型由val的类型决定; map[xxx]string返回[]string, map[xxx]int返回[]int
 * 示例：targetMap := map[string]int{}, valSlice := MapValSlice(targetMap).([]int)
 */
func MapValSlice(m interface{}) interface{} {
    typ, val := reflect.TypeOf(m), reflect.ValueOf(m)
    if typ.Kind() != reflect.Map {
        return nil
    }
    
    keys := val.MapKeys()
    slice := reflect.New(reflect.SliceOf(typ.Elem())).Elem()
    for _, key := range keys {
        val := val.MapIndex(key)
        slice.Set(reflect.Append(slice, val))
    }
    return slice.Interface()
}

/*
 * 深拷贝;
 * [指针类型字段, map, slice]全部新建, 确保修改 拷贝后对象的任意字段 不会影响 源对象对应字段
 * 未导出字段(小写字母开头)也可以拷贝
 */
func DeepCopy(item interface{}) (newItem interface{}, err error) {
    
    srcTyp, srcVal := reflect.TypeOf(item), reflect.ValueOf(item)
    if srcTyp.Kind() != reflect.Ptr {
        return nil, errors.New("is not pointer")
    }
    
    defer func() {
        if e := recover(); e != nil {
            err = fmt.Errorf("%v", e)
        }
    }()
    
    srcTyp, srcVal = srcTyp.Elem(), srcVal.Elem()
    newValPtr := reflect.New(srcTyp)
    newVal := newValPtr.Elem()
    
    for i := 0; i < srcTyp.NumField(); i++ {
        field, fieldVal := srcTyp.Field(i), srcVal.Field(i)
        if field.PkgPath == "" {
            newVal.Field(i).Set(deepCopyValue(field.Name, field.Type, &fieldVal))
        } else {
            forceSrcVal := reflect.NewAt(field.Type, unsafe.Pointer(fieldVal.UnsafeAddr())).Elem()
            forceNewVal := reflect.NewAt(field.Type, unsafe.Pointer(newVal.Field(i).UnsafeAddr())).Elem()
            forceNewVal.Set(deepCopyValue(field.Name, field.Type, &forceSrcVal))
        }
    }
    return newValPtr.Interface(), nil
}

func deepCopyValue(fieldName string, typ reflect.Type, val *reflect.Value) reflect.Value {
    switch typ.Kind() {
    
    case reflect.Ptr:
        if val.IsNil() {
            return *val
        }
        v := val.Elem()
        return deepCopyValue(fieldName, typ.Elem(), &v).Addr()
    
    case reflect.Interface:
        if val.IsNil() {
            return *val
        }
        v := val.Elem()
        return deepCopyValue(fieldName, v.Type(), &v).Addr()
    
    case reflect.String:
        // 可能是 string类型的别名，所以不要valueOf(string)
        // s := val.String()
        // return reflect.ValueOf(&s).Elem()
        stringPtr := reflect.New(typ)
        stringPtr.Elem().SetString(val.String())
        return stringPtr.Elem()
    
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        intPtr := reflect.New(typ)
        intPtr.Elem().SetInt(val.Int())
        return intPtr.Elem()
    
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        uIntPtr := reflect.New(typ)
        uIntPtr.Elem().SetUint(val.Uint())
        return uIntPtr.Elem()
    
    case reflect.Struct:
        st := val.Addr().Interface()
        if v, err := DeepCopy(st); err != nil {
            panic(err.Error())
        } else {
            return reflect.ValueOf(v).Elem()
        }
    
    case reflect.Map:
        newMap := reflect.MakeMap(typ)
        for _, srcKey := range val.MapKeys() {
            newKey := deepCopyValue(fieldName+".key", srcKey.Type(), &srcKey)
            srcVal := val.MapIndex(srcKey)
            newVal := deepCopyValue(fieldName+".val", srcVal.Type(), &srcVal)
            newMap.SetMapIndex(newKey, newVal)
        }
        newMapPtr := reflect.New(typ)
        newMapPtr.Elem().Set(newMap)
        return newMapPtr.Elem()
    
    case reflect.Slice:
        newSlice := reflect.MakeSlice(typ, val.Len(), val.Cap())
        for i := 0; i < val.Len(); i++ {
            srcVal := val.Index(i)
            newVal := deepCopyValue(fieldName+".item", srcVal.Type(), &srcVal)
            newSlice.Index(i).Set(newVal)
        }
        newSlicePtr := reflect.New(typ)
        newSlicePtr.Elem().Set(newSlice)
        return newSlicePtr.Elem()
    
    default:
        panic("unsupported field[" + fieldName + "] type: " + typ.Kind().String())
    }
}
