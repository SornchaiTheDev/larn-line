package utils

func Has(slice []string, item string) bool {
	isHas := false

	for _, _item := range slice {
		if _item == item {
			isHas = true
			break
		}
	}

	return isHas
}
