package expressions

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTokenizeToNumberMonths(t *testing.T) {
	input := "1 2 3 JUN-DEC 4 cmd"
	expected := &CronElements{
		Minute:   "1",
		Hour:     "2",
		DayMonth: "3",
		Month:    "6-12",
		DayWeek:  "4",
		Command:  "cmd",
	}

	ds, err := NewDefaultSyntax(input)
	assert.Nil(t, err)
	d, _ := ds.(*DefaultSyntax)

	err = d.tokenize()
	assert.Nil(t, err)
	actual := d.cronElements
	assert.Equal(t, expected, actual)

}

func TestStringToNumberNoop(t *testing.T) {
	input := "1 2 3 6-12 4 cmd"
	expected := &CronElements{
		Minute:   "1",
		Hour:     "2",
		DayMonth: "3",
		Month:    "6-12",
		DayWeek:  "4",
		Command:  "cmd",
	}
	ds, err := NewDefaultSyntax(input)
	assert.Nil(t, err)
	d, _ := ds.(*DefaultSyntax)

	err = d.tokenize()
	assert.Nil(t, err)
	actual := d.cronElements
	assert.Equal(t, expected, actual)
}

func TestStringToNumberDays(t *testing.T) {
	input := "1/15 8,9-10 * 6-12 SUN-WED cmd"
	expected := &CronElements{
		Minute:   "1/15",
		Hour:     "8,9-10",
		DayMonth: "*",
		Month:    "6-12",
		DayWeek:  "0-3",
		Command:  "cmd",
	}

	ds, err := NewDefaultSyntax(input)
	assert.Nil(t, err)
	d, _ := ds.(*DefaultSyntax)

	err = d.tokenize()
	assert.Nil(t, err)
	actual := d.cronElements
	assert.Equal(t, expected, actual)
}

func TestValidateTokens(t *testing.T) {
	input := "*/15 0 1-5 4 1-3/7 command"
	ds, err := NewDefaultSyntax(input)
	assert.Nil(t, err)
	d, _ := ds.(*DefaultSyntax)

	tcs := []struct {
		name  string // name of the test
		input string // string to validate
		error bool   // expected error, true if validation failed
	}{
		{"single int", "19", false},
		{"*", "*", false},
		{"interval", "19-26", false},
		{"step with *", "*/10", false},
		{"step with single number", "2/10", false},
		{"step with interval", "2-9/10", false},
		{"multi tokens", "2-9/10,4,3-10", false},
		{"invalid string", "a", true},
		{"invalid interval", "-4", true},
		{"invalid * in interval start", "*-4", true},
		{"invalid * in interval end", "4-*", true},
		{"same interval", "4-4", false},
		{"invalid step with inteval", "4/4-6", true},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			actual := d.validateTokens(tc.input)
			actualError := false
			if actual != nil {
				actualError = true
			}
			assert.Equal(t, actualError, tc.error)
		})
	}
}

func TestElements(t *testing.T) {
	input := "*/15 0 1-5 4 1-3/7 command"
	ds, err := NewDefaultSyntax(input)
	require.Nil(t, err)
	actual, err := ds.Elements()
	assert.Nil(t, err)
	assert.NotNil(t, actual)
}

func TestValidateExpressionWrongFields(t *testing.T) {
	input := "1 2 3 4 5"
	ds, err := NewDefaultSyntax(input)
	require.Nil(t, err)
	actual := ds.ValidateExpression(input)
	assert.NotNil(t, actual)
}

func TestValidateExpressionWrongTokens(t *testing.T) {
	input := "*/1-8 2 3 4 5 /bin/ls"
	ds, err := NewDefaultSyntax(input)
	require.Nil(t, err)
	actual := ds.ValidateExpression(input)
	assert.NotNil(t, actual)
}
