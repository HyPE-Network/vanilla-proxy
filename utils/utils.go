package utils

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

func Format(a []any) string {
	return strings.TrimSuffix(strings.TrimSuffix(fmt.Sprintln(a...), "\n"), "\n")
}

func GetPluralForm(num int, form1 string, form2 string, form3 string) string {
	n := int(math.Abs(float64(num))) % 100
	n1 := n % 10
	if n > 10 && n < 20 {
		return form3
	} else if n1 > 1 && n1 < 5 {
		return form2
	} else if n1 == 1 {
		return form1
	}
	return form3
}

func GetFullPluralForm(num int, form1 string, form2 string, form3 string) string {
	n := int(math.Abs(float64(num))) % 100
	n1 := n % 10
	if n > 10 && n < 20 {
		return strconv.Itoa(num) + " " + form3
	} else if n1 > 1 && n1 < 5 {
		return strconv.Itoa(num) + " " + form2
	} else if n1 == 1 {
		return strconv.Itoa(num) + " " + form1
	}
	return strconv.Itoa(num) + " " + form3
}

func GetTimestamp() int64 {
	return time.Now().Unix()
}

func GetMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func MillisToDate(t int64) string {
	return TimestampToDate(t / 1000)
}

func TimestampToDate(t int64) string {
	return FormatTime(time.Unix(t, 0))
}

func FormatTime(t time.Time) string {
	return t.Format("15:04:05 2006.01.02") // magic numbers https://go.dev/src/time/format.go
}
