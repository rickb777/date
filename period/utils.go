package period

func allInt(f func(i int) bool, items ...int) bool {
	for _, item := range items {
		if !f(item) {
			return false
		}
	}

	return true
}

func anyInt(f func(i int) bool, items ...int) bool {
	for _, item := range items {
		if f(item) {
			return true
		}
	}

	return false
}

func intEqualZero(i int) bool {
	return i == 0
}

func intGreaterZero(i int) bool {
	return i > 0
}

func intLessZero(i int) bool {
	return i < 0
}
