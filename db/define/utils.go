package define

// InString in字符串操作
func InString(arr []string) []string {
	data := make([]string, 0)
	for _, s := range arr {
		data = append(data, "'"+s+"'")
	}
	return data
}
