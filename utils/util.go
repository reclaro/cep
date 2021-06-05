package utils

import (
	"sort"
	"strings"
)

// RangeValues receives an int of two values and return all the values between interval[0] and interval[1] included
func RangeValues(interval []int) []int {
	resu := []int{}
	for i := interval[0]; i <= interval[1]; i++ {
		resu = append(resu, i)
	}
	return resu
}

// StringToNumber replace in an input string all the occurrences of the key defined in the mapper with
// the corresponding value
func StringToNumber(input string, mapper map[string]string) string {
	lower := strings.ToUpper(input)
	for str, num := range mapper {
		lower = strings.ReplaceAll(lower, str, num)
	}
	return lower
}

// RemoveDups remove the duplicates from an array of int
func RemoveDups(input []int) []int {
	tmpHash := make(map[int]bool)
	for _, v := range input {

		tmpHash[v] = true
	}
	noDupes := []int{}
	for k, _ := range tmpHash {
		noDupes = append(noDupes, k)
	}
	return noDupes
}

// SortedUniqueInts receives an array of int and return the array without the duplicates and
// in an ascending order
func SortedUniqueInts(input []int) []int {
	input = RemoveDups(input)
	sort.Ints(input)
	return input
}
