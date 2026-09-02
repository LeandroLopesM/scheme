package parser

func stringAround(rangeSize int, vec []rune, middle int) string {
	back := middle - rangeSize
	if middle-rangeSize < 0 {
		back = 0
	}

	front := middle + rangeSize
	if middle+rangeSize >= len(vec) {
		front = len(vec) - 1
	}

	return string(vec[back:front])
}
