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

// 求最小值
func IntMin(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// 求最大值
func IntMax(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// 返回最大值对应的序号
func Max(values []float64) (no int) {
	maxValue := values[0]
	no = 0
	for i := 1; i < len(values); i++ {
		if values[i] > maxValue {
			maxValue = values[i]
			no = i
		}
	}
	return
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

// const FLOAT_ZERO = 10e-20 // 定义接近0的极小值
const FLOAT_ZERO = math.SmallestNonzeroFloat64 * 10

// 浮点数相等比较
func IsEqual(value1 float64, value2 float64) bool {
	if math.Abs(value1-value2) < FLOAT_ZERO {
		return true
	}
	return false
}

// 浮点数大于等于
func IsBigEqual(value1 float64, value2 float64) bool {
	if value1 > value2 || math.Abs(value1-value2) < FLOAT_ZERO {
		return true
	}
	return false
}

// 浮点数小于等于
func IsSmallEqual(value1 float64, value2 float64) bool {
	if value1 > value2 || math.Abs(value1-value2) < FLOAT_ZERO {
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

type Int64s []int64

//Len()
func (s Int64s) Len() int {
	return len(s)
}

//Less():成绩将有低到高排序
func (s Int64s) Less(i, j int) bool {
	return s[i] < s[j]
}

//Swap()
func (s Int64s) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
