package functions

func FilterInPlace(
	arr []interface{},
	fn func(t int) bool,
) int {
	deleteIndex := -1
	for t:=0; t<len(arr); t++ {
		if fn(t) {
			if deleteIndex != -1 {
				arr[deleteIndex], arr[t] = arr[t], arr[deleteIndex]
				deleteIndex+=1
			}
		} else {
			if deleteIndex == -1 {
				deleteIndex = t
			}
		}
	}

	if deleteIndex == -1 {
		return len(arr)
	}
	return deleteIndex
}