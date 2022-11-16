package router

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// GetBody 获取body数据
func GetBody(r *http.Request) string {
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	_ = r.Body.Close()

	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	return string(bodyBytes)
}

// ReadJSON 读取 Json Body
func ReadJSON(r *http.Request, data interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	//	保存数据到结构体
	if err = json.Unmarshal(body, data); err != nil {
		return err
	}

	return nil
}

// FormatJSON JSON格式数据
type FormatJSON struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// IndexData 列表数据
type IndexData struct {
	Items interface{} `json:"items"`
	Count int64       `json:"count"`
}

// SuccessJSON 返回成功数据
func SuccessJSON(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(http.StatusOK)
	resp := &FormatJSON{
		Code: 0,
		Msg:  "",
		Data: data,
	}
	s, _ := json.Marshal(resp)
	_, _ = w.Write(s)
}

// ErrorJSON 返回错误数据
func ErrorJSON(w http.ResponseWriter, msg string, code int) {
	w.WriteHeader(http.StatusOK)
	resp := &FormatJSON{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
	s, _ := json.Marshal(resp)
	_, _ = w.Write(s)
}
