package utils

import (
	"fmt"
	"strings"
)

// SliceStringIndexOf 数组中是否包含字符串
func SliceStringIndexOf(seq string, slice []string) int {
	for index, val := range slice {
		if seq == val {
			return index
		}
	}
	return -1
}

// SliceInt64IndexOf 数组中是否包含整数
func SliceInt64IndexOf(seq int64, slice []int64) int {
	for index, val := range slice {
		if seq == val {
			return index
		}
	}
	return -1
}

// MapBuildQuery map转请求参数
func MapBuildQuery(data map[string]interface{}) string {
	var s []string
	for k, v := range data {
		s = append(s, fmt.Sprintf("%v=%v", k, v))
	}
	return strings.Join(s, "&")
}

// GetMapKeys 获取Map key所有值
func GetMapKeys(oldMap map[string]string) []string {
	keys := make([]string, 0, len(oldMap))
	for k, _ := range oldMap {
		keys = append(keys, k)
	}
	return keys
}

// GetMapValues 获取Map val所有值
func GetMapValues(oldMap map[string]string) []string {
	values := make([]string, 0, len(oldMap))
	for _, v := range oldMap {
		values = append(values, v)
	}

	return values
}
