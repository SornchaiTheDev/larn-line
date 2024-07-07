package utils

func IndexOf(slice string, item rune) int {
	idx := -1
	for index, char := range slice {
		if char == item {
			idx = index
			break
		}
	}

	return idx
}
