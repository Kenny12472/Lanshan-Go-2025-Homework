func countarray(arr []int) map[int]int {
	count := make(map[int]int)
	for i := 0; i < len(arr); i++ {
		value, exists := count[arr[i]]
		if exists {
			count[arr[i]] = value + 1
		} else {
			count[arr[i]] = 1
		}
	}
	return count
}

