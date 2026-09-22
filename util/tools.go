package util

import (
	"math"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// GetCurrentFilename returns the name, without extension, of the Go file that
// calls it.
func GetCurrentFilename() string {
	_, callerFile, _, _ := runtime.Caller(1)
	baseFile := filepath.Base(callerFile)
	return strings.TrimSuffix(baseFile, filepath.Ext(baseFile))
}

type byteUnit struct {
	threshold float64
	size      float64
	symbol    string
}

var byteUnits = []byteUnit{
	{threshold: 1e15, size: 1 << 50, symbol: "PB"},
	{threshold: 1e12, size: 1 << 40, symbol: "TB"},
	{threshold: 1e9, size: 1 << 30, symbol: "GB"},
	{threshold: 1e6, size: 1 << 20, symbol: "MB"},
	{threshold: 1e3, size: 1 << 10, symbol: "KB"},
}

// ByteFormat renders a byte count in the largest fitting unit, rounded up to
// the given number of decimal places.
// Adapted from https://www.socketloop.com/tutorials/golang-byte-format-example
func ByteFormat(byteCount int64, precision int) string {
	if precision <= 0 {
		precision = 1
	}
	count := float64(byteCount)
	for _, unit := range byteUnits {
		if count >= unit.threshold {
			return formatFloat(roundUp(count/unit.size, precision)) + " " + unit.symbol
		}
	}
	return formatFloat(count) + " B"
}

func formatFloat(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func roundUp(value float64, places int) float64 {
	scale := math.Pow(10, float64(places))
	return math.Ceil(value*scale) / scale
}

// ReplaceLast replaces the last occurrence of old in s with replacement.
func ReplaceLast(s, old, replacement string) string {
	return reverse(strings.Replace(reverse(s), reverse(old), reverse(replacement), 1))
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
