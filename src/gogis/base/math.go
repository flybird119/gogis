package base

import (
	"math"
	"strings"
	"time"
)

// 求绝对值
func Abs(x int) int {
	if x >= 0 {
		return x
	}
	return -x
}

func IntMin(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func IntMax(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func IsMatchBool(value1 bool, op string, value2 bool) bool {
	switch op {
	case "=":
		return value1 == value2
	case "!=":
		return value1 != value2
	}
	return false
}

func IsMatchInt(value1 int, op string, value2 int) bool {
	switch op {
	case "=":
		return value1 == value2
	case "!=":
		return value1 != value2
	case ">":
		return value1 > value2
	case "<":
		return value1 < value2
	case ">=":
		return value1 >= value2
	case "<=":
		return value1 <= value2
	}
	return false
}

const FLOAT_ZERO = 10e-10 // 定义接近0的极小值
func IsEqual(value1 float64, value2 float64) bool {
	if math.Abs(value1-value2) < FLOAT_ZERO {
		return true
	}
	return false
}

func IsMatchFloat(value1 float64, op string, value2 float64) bool {
	switch op {
	case "=":
		return IsEqual(value1, value2)
	case "!=":
		return !IsEqual(value1, value2)
	case ">":
		return value1 > value2
	case "<":
		return value1 < value2
	case ">=":
		return value1 >= value2
	case "<=":
		return value1 <= value2
	}
	return false
}

func IsMatchString(value1 string, op string, value2 string) bool {
	switch op {
	case "=":
		return strings.EqualFold(value1, value2)
	case "!=":
		return !strings.EqualFold(value1, value2)
		// 感觉没啥用，先封起来
		// case ">":
		// 	return value1 > value2
		// case "<":
		// 	return value1 < value2
		// case ">=":
		// 	return value1 >= value2
		// case "<=":
		// 	return value1 <= value2
	}
	return false
}

func IsMatchTime(value1 time.Time, op string, value2 time.Time) bool {
	switch op {
	case "=":
		return value1.Equal(value2)
	case "!=":
		return !value1.Equal(value2)
	case ">":
		return value1.After(value2)
	case "<":
		return value1.Before(value2)
	case ">=":
		return value1.Equal(value2) || value1.After(value2)
	case "<=":
		return value1.Before(value2) || value1.After(value2)
	}
	return false
}

// 计算两点距离的平方
func DistanceSquare(x0, y0, x1, y1 float64) float64 {
	return math.Pow((x0-x1), 2) + math.Pow((y0-y1), 2)
}
