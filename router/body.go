package router

import (
	"encoding/json"
	"net/http"
)

// FormatJSON JSON格式数据
type FormatJSON struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
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
