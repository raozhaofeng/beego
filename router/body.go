package router

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// RangeTimeParam 时间范围参数
type RangeTimeParam struct {
	Form string `json:"form"` //	开始时间
	To   string `json:"to"`   //	结束时间
}

// Pagination 分页
type Pagination struct {
	SortBy      string `json:"sortBy"`      //	排序字段
	Descending  bool   `json:"descending"`  //	是否降序排序
	Page        int64  `json:"page"`        //	当前页数
	RowsPerPage int64  `json:"rowsPerPage"` //	每页显示条数
}

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
