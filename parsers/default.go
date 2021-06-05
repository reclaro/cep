package parsers

import (
	"fmt"
	"github.com/reclaro/cep/validators"
)

// test case https://crontab.guru/#1-19/5,*/10,50-59,*_3_8_JAN_TUE
type Parser interface {
	Minute() []int
	Hour() []int
	DayOfTheMonth() []int
	Month() []string
	DayOfTheWeek() []string
	Command() string
}

type DefaultParser struct {
	selectors        []string
	minsValues       []int
	hoursValues      []int
	dayOfMonthValues []int
	monthInt         []int
	monthString      map[string]int
	dayOfWeekInt     []int
	dayOfWeekString  map[string]int
	validator        validators.Validator
	inputLine        string
}

func NewDefaultParser(syntaxValidator validators.Validator, input string) Parser {
	dp := &DefaultParser{
		selectors:        []string{"-", "/", "*", ","},
		minsValues:       []int{0, 59},
		hoursValues:      []int{0, 23},
		dayOfMonthValues: []int{1, 31},
		monthInt:         []int{1, 31},
		monthString: map[string]int{"JAN": 1,
			"FEB": 2,
			"MAR": 3,
			"APR": 4,
			"MAY": 5,
			"JUN": 6,
			"JUL": 7,
			"AUG": 8,
			"SEP": 9,
			"OCT": 10,
			"NOV": 11,
			"DEC": 12},
		dayOfWeekInt:    []int{0, 6},
		dayOfWeekString: map[string]int{"SUN": 0, "MON": 1, "TUE": 2, "WED": 3, "THU": 4, "FRI": 5, "SAT": 6},
		validator:       syntaxValidator,
		inputLine:       input,
	}

	dp.validator.Tokenize(inputLine)
	return dp
}

//func (dp *DefaultParser)
