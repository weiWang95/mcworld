package util

func FloorFloat(d float32) int64 {
	// 2.1 => 2
	if d >= 0 {
		return int64(d)
	}

	// -1.0 => -1
	if d-float32(int64(d)) == 0 {
		return int64(d)
	}

	// -1.2 => -2
	return int64(d - 1)
}

func MinInt64(a, b int64) int64 {
	if a <= b {
		return a
	}

	return b
}

func MaxInt64(a, b int64) int64 {
	if a >= b {
		return a
	}

	return b
}
