package utils

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
