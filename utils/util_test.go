package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringToNumber(t *testing.T) {
	mapper := map[string]string{"MON": "1", "TUE": "2", "WED": "3", "THU": "4"}
	input := "TUE-WED"
	expected := "2-3"
	actual := StringToNumber(input, mapper)
	assert.Equal(t, expected, actual)
}

func TestStringToNumberCaseInsesitive(t *testing.T) {
	mapper := map[string]string{"MON": "1", "TUE": "2", "WED": "3", "THU": "4"}
	input := "tue-WED"
	expected := "2-3"
	actual := StringToNumber(input, mapper)
	assert.Equal(t, expected, actual)
}

func TestStringToNumberNoop(t *testing.T) {
	mapper := map[string]string{"MON": "1", "TUE": "2", "WED": "3", "THU": "4"}
	input := "THIS * 2,3 A/ -STRING"
	expected := input
	actual := StringToNumber(input, mapper)
	assert.Equal(t, expected, actual)
}

func TestRemoveDups(t *testing.T) {
	input := []int{0, 1, 2, 3, 4, 3, 0}
	expected := []int{0, 1, 2, 3, 4}
	actual := RemoveDups(input)
	assert.Equal(t, len(expected), len(actual))
}

func TestRemoveDupsNoop(t *testing.T) {
	input := []int{0, 1, 2, 3}
	actual := RemoveDups(input)
	assert.Equal(t, len(input), len(actual))
}

func TestSortedUniqueInts(t *testing.T) {
	input := []int{0, 1, 2, 3, 17, 8, 21, 1}
	expected := []int{0, 1, 2, 3, 8, 17, 21}
	actual := SortedUniqueInts(input)
	assert.Equal(t, expected, actual)
}
