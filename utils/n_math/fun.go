package n_math

import "math/rand"

func Min64(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max64(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func GetRandInt32() int32 {
	return rand.Int31()
}

func getRandInt8() int8  {
	return int8(rand.Int31n(256))
}