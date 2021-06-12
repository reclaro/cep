package parsers

import (
	"strings"
	"testing"

	"github.com/reclaro/cep/expressions"
	"github.com/reclaro/cep/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func defaultParserWithDefaultHolder(t *testing.T) *DefaultParser {
	input := "*/15 0 1,15 * 1-5 /usr/bin/find"
	holder, e := expressions.NewDefaultSyntax(input)
	assert.Nil(t, e)
	p, e := NewDefaultParser(holder)
	dp, ok := p.(*DefaultParser)
	require.True(t, ok)
	assert.Nil(t, e)
	return dp
}

func defaultParserWithDefaultHolderWithString(t *testing.T, input string) *DefaultParser {

	holder, e := expressions.NewDefaultSyntax(input)
	assert.Nil(t, e)
	p, e := NewDefaultParser(holder)
	dp, ok := p.(*DefaultParser)
	require.True(t, ok)
	assert.Nil(t, e)
	return dp
}

func TestParseAllValues(t *testing.T) {
	dp := defaultParserWithDefaultHolder(t)
	allValues := []int{0, 59}
	expected := utils.RangeValues(allValues)
	actual, err := dp.parse("18,*", allValues)
	assert.Nil(t, err)

	assert.Equal(t, expected, actual)
}

func TestManageStepsNoInterval(t *testing.T) {
	dp := defaultParserWithDefaultHolder(t)
	allowedValues := []int{0, 59}
	input := "4/10"
	expected := []int{4, 14, 24, 34, 44, 54}
	actual, err := dp.manageSteps(strings.Split(input, "/"), allowedValues)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}
func TestManageStepsAllValues(t *testing.T) {
	dp := defaultParserWithDefaultHolder(t)
	allowedValues := []int{0, 59}
	input := "*/10"
	expected := []int{0, 10, 20, 30, 40, 50}
	actual, err := dp.manageSteps(strings.Split(input, "/"), allowedValues)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestManageStepsWithInterval(t *testing.T) {
	dp := defaultParserWithDefaultHolder(t)
	allowedValues := []int{0, 59}
	input := "1-19/10"
	expected := []int{1, 11}
	actual, err := dp.manageSteps(strings.Split(input, "/"), allowedValues)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestManageStepsInvalidStep(t *testing.T) {
	dp := defaultParserWithDefaultHolder(t)
	allowedValues := []int{0, 59}
	input := "1-19/*"
	_, err := dp.manageSteps(strings.Split(input, "/"), allowedValues)
	assert.NotNil(t, err)
}

func TestManageIntervals(t *testing.T) {
	dp := defaultParserWithDefaultHolder(t)
	expected := []int{2, 45}
	input := []string{"2", "45"}
	actual, err := dp.manageIntervals(input)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestManageIntervalsInvalidString(t *testing.T) {
	dp := defaultParserWithDefaultHolder(t)
	input := []string{"A", "45"}
	_, err := dp.manageIntervals(input)
	assert.NotNil(t, err)
}

func TestManageIntervalsInvalidInterval(t *testing.T) {
	dp := defaultParserWithDefaultHolder(t)
	input := []string{"3", "1"}
	_, err := dp.manageIntervals(input)
	assert.NotNil(t, err)
}

func TestManageIntervalsStartEqEnd(t *testing.T) {
	dp := defaultParserWithDefaultHolder(t)

	input := []string{"3", "3"}
	_, err := dp.manageIntervals(input)
	assert.Nil(t, err)
}

func TestParseInterval(t *testing.T) {
	dp := defaultParserWithDefaultHolder(t)
	input := "*/15,3-11"

	actual, err := dp.parse(input, dp.minsValues)
	assert.Nil(t, err)
	expected := []int{0, 15, 30, 45, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	assert.Equal(t, expected, actual)
}
func TestParseIntervalStepSingleValue(t *testing.T) {
	dp := defaultParserWithDefaultHolder(t)
	input := "*/15,3-11,59"

	actual, err := dp.parse(input, dp.minsValues)
	assert.Nil(t, err)
	expected := []int{0, 15, 30, 45, 3, 4, 5, 6, 7, 8, 9, 10, 11, 59}
	assert.Equal(t, expected, actual)
}

func TestParseSingleValue(t *testing.T) {
	dp := defaultParserWithDefaultHolder(t)
	input := "0"

	actual, err := dp.parse(input, dp.minsValues)
	assert.Nil(t, err)
	expected := []int{0}
	assert.Equal(t, expected, actual)
}

func TestParseOnlyInterval(t *testing.T) {
	dp := defaultParserWithDefaultHolder(t)
	input := "1-9"

	actual, err := dp.parse(input, dp.minsValues)
	assert.Nil(t, err)
	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	assert.Equal(t, expected, actual)
}

func TestParseMaySeparateValues(t *testing.T) {
	dp := defaultParserWithDefaultHolder(t)
	input := "1,9,56"

	actual, err := dp.parse(input, dp.minsValues)
	assert.Nil(t, err)
	expected := []int{1, 9, 56}
	assert.Equal(t, expected, actual)
}

func TestParseComplexValuesAll(t *testing.T) {
	dp := defaultParserWithDefaultHolder(t)
	// The * should make this returns all values
	input := "1-19/5,*/10,50-59,*"

	actual, err := dp.parse(input, dp.minsValues)
	assert.Nil(t, err)
	expected := utils.RangeValues(dp.minsValues)
	assert.Equal(t, expected, actual)
}

func TestParseComplexValues(t *testing.T) {
	dp := defaultParserWithDefaultHolder(t)
	input := "1-19/5,*/10,50-59"

	actual, err := dp.parse(input, dp.minsValues)
	require.Nil(t, err)
	expected := []int{0, 1, 6, 10, 11, 16, 20, 30, 40}
	expected = append(expected, utils.RangeValues([]int{50, 59})...)
	actual = utils.SortedUniqueInts(actual)
	assert.Equal(t, len(expected), len(actual))
	assert.Equal(t, expected, actual)
}

func TestInAllowedValues(t *testing.T) {
	dp := defaultParserWithDefaultHolder(t)
	input := []int{1, 19}
	assert.True(t, dp.inAllowedValues(input, dp.minsValues))
}

func TestInAllowedValuesOutOfRange(t *testing.T) {
	dp := defaultParserWithDefaultHolder(t)
	input := []int{1, 60}
	// error on the max value
	assert.False(t, dp.inAllowedValues(input, []int{0, 59}))
	// error on the min valus
	input = []int{-1, 11}
	assert.False(t, dp.inAllowedValues(input, []int{0, 12}))
}

func TestMinutes(t *testing.T) {
	input := "*/15 0 1,15 * 1-5 /usr/bin/find"
	ex, err := expressions.NewDefaultSyntax(input)
	require.Nil(t, err)
	p, e := NewDefaultParser(ex)
	require.Nil(t, e)
	expected := []int{0, 15, 30, 45}
	actual, err := p.Minutes()
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestHours(t *testing.T) {
	input := "*/15 0 1,15 * 1-5 /usr/bin/find"
	ex, err := expressions.NewDefaultSyntax(input)
	require.Nil(t, err)
	p, e := NewDefaultParser(ex)
	require.Nil(t, e)
	expected := []int{0}
	actual, err := p.Hours()
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestDayOfMonth(t *testing.T) {

	input := "*/15 0 1,15 * 1-5 /usr/bin/find"
	ex, err := expressions.NewDefaultSyntax(input)
	require.Nil(t, err)
	p, e := NewDefaultParser(ex)
	require.Nil(t, e)
	expected := []int{1, 15}
	actual, err := p.DaysOfTheMonth()
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestMonth(t *testing.T) {
	input := "*/15 0 1,15 * 1-5 /usr/bin/find"
	ex, err := expressions.NewDefaultSyntax(input)
	require.Nil(t, err)
	p, e := NewDefaultParser(ex)
	require.Nil(t, e)
	expected := utils.RangeValues([]int{1, 12})
	actual, err := p.Months()
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestDayOfWeek(t *testing.T) {

	input := "*/15 0 1,15 * 1-5 /usr/bin/find"
	ex, err := expressions.NewDefaultSyntax(input)
	require.Nil(t, err)
	p, e := NewDefaultParser(ex)
	require.Nil(t, e)
	expected := utils.RangeValues([]int{1, 5})
	actual, err := p.DaysOfTheWeek()
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestCommand(t *testing.T) {
	input := "*/15 0 1,15 * 1-5 /usr/bin/find"
	ex, err := expressions.NewDefaultSyntax(input)
	require.Nil(t, err)
	p, e := NewDefaultParser(ex)
	require.Nil(t, e)
	expected := "/usr/bin/find"
	actual, err := p.Command()
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestOutOfRangeValues(t *testing.T) {
	brokenInput := "60 25 32 0 7 /bin/ls"
	ex, err := expressions.NewDefaultSyntax(brokenInput)
	require.Nil(t, err)
	p, e := NewDefaultParser(ex)
	require.Nil(t, e)

	_, errRet := p.Minutes()
	assert.NotNil(t, errRet)

	_, errRet = p.Hours()
	assert.NotNil(t, errRet)

	_, errRet = p.DaysOfTheMonth()
	assert.NotNil(t, errRet)

	_, errRet = p.DaysOfTheWeek()
	assert.NotNil(t, errRet)

	_, errRet = p.Months()
	assert.NotNil(t, errRet)

}

func TestResultsValid(t *testing.T) {
	dp := defaultParserWithDefaultHolder(t)
	expected, err := dp.Results()
	require.Nil(t, err)
	assert.NotNil(t, expected)
}

func TestResultsInvalid(t *testing.T) {
	// input := "*/15 0 1,15 * 99 /usr/bin/find"
	// dp := defaultParserWithDefaultHolderWithString(t, input)
	// expected, err := dp.Results()
	// assert.NotNil(t, err)
	// assert.Nil(t, expected)

	tcs := []struct {
		name  string
		input string
	}{
		{"Invalid mins", "99 * * * * /usr/bin/find"},
		{"Invalid hours", "* 99 * * * /usr/bin/find"},
		{"Invalid days of month", "* * 99 * * /usr/bin/find"},
		{"Invalid monts", "* * * 99 * /usr/bin/find"},
		{"Invalid days of week", "* * * * 99 /usr/bin/find"},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			dp := defaultParserWithDefaultHolderWithString(t, tc.input)
			expected, err := dp.Results()
			assert.NotNil(t, err)
			assert.Nil(t, expected)

		})
	}
}
