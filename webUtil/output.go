/*
Created by
User: junjunzhou
Time: 2017/9/21 17:49
*/

package webUtil

import (
	"bytes"
	"encoding/json"
	"github.com/ZhouJunjun/goLib/log4j"
	"github.com/ZhouJunjun/goLib/util/logBuffer"
	"html/template"
	"net/http"
	"regexp"
	"strings"
)

const (
	CONTENT_TYPE_JSON_UTF8       = "application/json;charset=UTF-8"
	CONTENT_TYPE_JAVASCRIPT_UTF8 = "application/javascript;charset=UTF-8"
	CONTENT_TYPE_HTML_UTF8       = "text/html;charset=UTF-8"
	CONTENT_TYPE_PLAIN_UTF8      = "text/plain;charset=UTF-8"
)

// 输出标准json格式的数据
func Output(writer http.ResponseWriter, status int, message string, data interface{}) {
	OutputAjax(writer, &AjaxResponse{Status: status, Message: message, Data: data})
}

// 输出标准javascript格式的数据
func OutputScript(writer http.ResponseWriter, callback string, status int, message string, data interface{}) {
	OutputScriptAjax(writer, callback, &AjaxResponse{Status: status, Message: message, Data: data})
}

// 输出标准html格式的数据
func OutputHtml(writer http.ResponseWriter, callback string, status int, message string, data interface{}) {
	OutputHtmlAjax(writer, callback, &AjaxResponse{Status: status, Message: message, Data: data})
}

//---------------------------------------------------------------------------------------------------------------------

// 输出标准json格式的数据
func OutputAndLog(log *logBuffer.Buffer, writer http.ResponseWriter, status int, message string, data interface{}) {
	ajax := &AjaxResponse{Status: status, Message: message, Data: data}
	if json, err := outputAjax(writer, ajax); err != nil {
		log.SetError(true).SetPrintStack(true).AppendFormat("; output error, data=%+v, err=%+v", ajax, err)
	} else {
		log.AppendString("; output=").Append(json)
	}
}

// 输出标准javascript格式的数据
func OutputScriptAndLog(log *logBuffer.Buffer, writer http.ResponseWriter, callback string, status int, message string, data interface{}) {
	if c := GetCallback(callback); c != "" {
		ajax := &AjaxResponse{Status: status, Message: message, Data: data}
		if _, json, err := outputScriptAjax(writer, callback, ajax); err != nil {
			log.SetError(true).SetPrintStack(true).AppendFormat("; output error, data=%+v, err=%+v", ajax, err)
		} else {
			log.AppendString("; output=").Append(json)
		}
	} else {
		OutputAndLog(log, writer, status, message, data)
	}
}

// 输出标准html格式的数据
func OutputHtmlAndLog(log *logBuffer.Buffer, writer http.ResponseWriter, callback string, status int, message string, data interface{}) {
	if c := GetCallback(callback); c != "" {
		ajax := &AjaxResponse{Status: status, Message: message, Data: data}
		if _, json, err := outputHtmlAjax(writer, callback, ajax); err != nil {
			log.SetError(true).SetPrintStack(true).AppendFormat("; output error, data=%+v, err=%+v", ajax, err)
		} else {
			log.AppendString("; output=").Append(json)
		}
	} else {
		OutputAndLog(log, writer, status, message, data)
	}
}

// --------------------------------------------------------------------------------------------------------------------

func OutputAjax(writer http.ResponseWriter, ajax Ajax) {
	if _, err := outputAjax(writer, ajax); err != nil {
		_ = log4j.ErrorStack("output error, data=%+v, err=%+v", ajax, err)
	}
}

func OutputScriptAjax(writer http.ResponseWriter, callback string, ajax Ajax) {
	if c := GetCallback(callback); c == "" {
		if _, err := outputAjax(writer, ajax); err != nil {
			_ = log4j.ErrorStack("output error, data=%+v, err=%+v", ajax, err)
		}
	} else {
		if _, _, err := outputScriptAjax(writer, c, ajax); err != nil {
			_ = log4j.ErrorStack("output error, data=%+v, err=%+v", ajax, err)
		}
	}
}

func OutputHtmlAjax(writer http.ResponseWriter, callback string, ajax Ajax) {
	if c := GetCallback(callback); c == "" {
		if _, err := outputAjax(writer, ajax); err != nil {
			_ = log4j.ErrorStack("output error, data=%+v, err=%+v", ajax, err)
		}
	} else {
		if _, _, err := outputHtmlAjax(writer, c, ajax); err != nil {
			_ = log4j.ErrorStack("output error, data=%+v, err=%+v", ajax, err)
		}
	}
}

// --------------------------------------------------------------------------------------------------------------------

func OutputBytes(writer http.ResponseWriter, json []byte, contentType string) {
	noCache(writer)
	writer.Header().Set("Content-Type", contentType)

	if _, err := writer.Write(json); err != nil {
		_ = log4j.ErrorStack("json encode error, data=%s, err=%+v", string(json), err)
	}
}

func OutputBytesAndLog(log *logBuffer.Buffer, writer http.ResponseWriter, json []byte, contentType string) {
	noCache(writer)
	writer.Header().Set("Content-Type", contentType)

	if _, err := writer.Write(json); err != nil {
		log.SetError(true).SetPrintStack(true).AppendFormat("; output error, data=%s, err=%+v", string(json), err)
	} else {
		log.AppendString("; output=").Append(json)
	}
}

// 后台输出页面模板
func OutPutTemplate(writer http.ResponseWriter, dir string, data map[string]interface{}, page ...string) {

	writer.Header().Set("Content-Type", "text/html;charset=UTF-8")
	writer.Header().Set("Pragma", "No-cache")
	writer.Header().Set("Cache-Control", "no-cache")

	if data == nil {
		data = map[string]interface{}{}
	}

	size := len(page)
	files := make([]string, size)
	for i, file := range page {
		if strings.HasPrefix(file, "/") {
			file = file[1:]
		}
		files[i] = dir + "/resources/" + file
	}

	if t3, err := template.ParseFiles(files...); err != nil {
		_ = log4j.ErrorStack("parseFiles error, %s", err.Error())
	} else {
		if err = t3.Execute(writer, data); err != nil {
			_ = log4j.ErrorStack("template execute error, %s", err.Error())
		}
	}
}

// 后台输出页面模板
func OutPutTemplateWithFunc(writer http.ResponseWriter, request *http.Request, dir string, funcMap template.FuncMap, data map[string]interface{}, page ...string) {

	writer.Header().Set("Content-Type", "text/html;charset=UTF-8")
	writer.Header().Set("Pragma", "No-cache")
	writer.Header().Set("Cache-Control", "no-cache")

	if data == nil {
		data = map[string]interface{}{}
	}

	size := len(page)
	files := make([]string, size)
	for i, file := range page {
		if strings.HasPrefix(file, "/") {
			file = file[1:]
		}
		files[i] = dir + "/resources/" + file
	}

	tmp := strings.Split(files[0], "/")
	name := tmp[len(tmp)-1]

	var err error = nil
	t := template.New(name).Funcs(funcMap)
	if t, err = t.ParseFiles(files...); err != nil {
		_ = log4j.ErrorStack("parseFiles error, %s", err.Error())
	} else if err := t.Execute(writer, data); err != nil {
		_ = log4j.ErrorStack("template execute error, %s", err.Error())
	}
}

// 华丽的分割线_______________________________________________________________________________________________________

// 输出字符串
func outputAjax(writer http.ResponseWriter, ajax Ajax) (json []byte, e error) {
	noCache(writer)
	writer.Header().Set("Content-Type", CONTENT_TYPE_JSON_UTF8)

	if js, err := ajax.Json(); err != nil {
		writer.WriteHeader(500)
		return nil, err
	} else if _, err = writer.Write(js); err != nil {
		writer.WriteHeader(500)
		return nil, err
	} else {
		return js, nil
	}
}

// 输出字符串
func outputScriptAjax(writer http.ResponseWriter, callback string, ajax Ajax) (script, json []byte, e error) {
	noCache(writer)
	writer.Header().Set("Content-Type", CONTENT_TYPE_JAVASCRIPT_UTF8)

	if sc, js, err := ajax.Script(callback); err != nil {
		writer.WriteHeader(500)
		return nil, nil, err

	} else if _, err = writer.Write(sc); err != nil {
		writer.WriteHeader(500)
		return nil, nil, err

	} else {
		return sc, js, nil
	}
}

// 输出字符串
func outputHtmlAjax(writer http.ResponseWriter, callback string, ajax Ajax) (html, json []byte, e error) {
	noCache(writer)
	writer.Header().Set("Content-Type", CONTENT_TYPE_HTML_UTF8)

	if ht, js, err := ajax.Html(callback); err != nil {
		writer.WriteHeader(500)
		return nil, nil, err

	} else if _, err = writer.Write(ht); err != nil {
		writer.WriteHeader(500)
		return nil, nil, err

	} else {
		return ht, js, nil
	}
}

// 设置无缓存
func noCache(writer http.ResponseWriter) {
	writer.Header().Set("Pragma", "No-cache")
	writer.Header().Set("Cache-Control", "no-cache")
}

func GetCallback(callback string) string {
	if callback != "" {
		ok, _ := regexp.MatchString("^[a-zA-Z0-9_]{1,64}$", callback)
		if ok {
			return callback
		}
	}
	return ""
}

type Ajax interface {
	Json() (js []byte, err error)
	Html(callback string) (html, js []byte, err error)
	Script(callback string) (script, js []byte, err error)
}

var (
	format_html_prefix = []byte(`<script>try {document.domain="56.com";} catch(e){};parent.`)
	format_html_suffix = []byte(`</script>`)
)

// 标准输出结构体
type AjaxResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewAjaxResponse(status int, msg string, data interface{}) *AjaxResponse {
	return &AjaxResponse{Status: status, Message: msg, Data: data}
}

func (p *AjaxResponse) Json() ([]byte, error) {
	if bs, err := json.Marshal(p); err != nil {
		return nil, err
	} else {
		return bs, nil
	}
}

func (p *AjaxResponse) Html(callback string) ([]byte, []byte, error) {

	bs, err := json.Marshal(p)
	if err != nil {
		return nil, nil, err
	}

	bf := bytes.NewBuffer(format_html_prefix)
	bf.WriteString(callback)
	bf.WriteByte('(')
	bf.Write(bs)
	bf.WriteByte(')')
	bf.Write(format_html_suffix)
	return bf.Bytes(), bs, nil
}

func (p *AjaxResponse) Script(callback string) ([]byte, []byte, error) {

	bs, err := json.Marshal(p)
	if err != nil {
		return nil, nil, err
	}

	bf := bytes.NewBufferString(callback)
	bf.WriteByte('(')
	bf.Write(bs)
	bf.WriteByte(')')
	return bf.Bytes(), bs, nil
}
