// Package view содержит логику для работы с пользовательским интерфейсом.
package view

func indexSwitch(current, count int) int {
	if current < count-1 {
		current++
	} else {
		current = 0
	}

	return current
}

func indexNext(current, count int) int {
	if current < count-1 {
		current++
	}

	return current
}

func indexPrev(current int) int {
	if current > 0 {
		current--
	}

	return current
}

func indexPrevWithCustomLimit(current, minVal int) int {
	if current > minVal {
		current--
	}

	return current
}
