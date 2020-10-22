package binsrch

func binsearch(srch int, sorted []int) int {
	keepLeft := func(i int) bool { return srch <= sorted[i] }
	idx := search(len(sorted), keepLeft)
	if idx >= len(sorted) || sorted[idx] != srch {
		return -1
	}

	return idx
}

func search(arySize int, keepLeft func(int) bool) int {
	idx, end := 0, arySize
	for idx < end {
		// divide by 2
		middle := int(uint(idx+end) >> 1)
		if keepLeft(middle) {
			// keep left part instead
			end = middle
		} else {
			// search is bigger than ary value, continue search at the left part
			idx = middle + 1
		}
	}
	return idx
}

func binsearch2(srch int, sorted []int) int {
	return recursiveBinSearch(sorted, srch, 0)
}

func recursiveBinSearch(sorted []int, srch, idx int) int {
	if len(sorted) == 0 {
		return -1
	}

	if len(sorted) == 1 {
		if sorted[0] != srch {
			return -1
		}
		return idx
	}

	middle := len(sorted) >> 1
	if srch == sorted[middle] {
		return idx + middle
	}

	if srch < sorted[middle] {
		sorted = sorted[0:middle]
	} else {
		sorted = sorted[middle:]
		idx = idx + middle
	}
	return recursiveBinSearch(sorted, srch, idx)
}
