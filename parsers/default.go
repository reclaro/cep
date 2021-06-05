package parsers

import (
	"errors"
	"fmt"
	"github.com/reclaro/cep/expressions"
	"github.com/reclaro/cep/utils"
	"strconv"
	"strings"
)

// CronResults contains the results of the parsing of a cron expression
type CronResults struct {
	Minute   []int
	Hour     []int
	DayMonth []int
	Month    []int
	DayWeek  []int
	Command  string
}

/*
   Parser define the methods to get the values of every single cron field and the method
   to get the CronResults
*/
type Parser interface {
	Minutes() ([]int, error)
	Hours() ([]int, error)
	DaysOfTheMonth() ([]int, error)
	Months() ([]int, error)
	DaysOfTheWeek() ([]int, error)
	Command() (string, error)
	Results() *CronResults
}

/* DefaultParser implements the Parser interface. The default parser follows the following rules
   minutes:  allowedValues 0-59
   hours:  allowed values 0-23
   days of the month: allowed values 1-31
   month: allowed Values 1-12
   day of the week: allowed Values 0-6. Sunday is the first day and it is reported as day 0
   In the DefaultParser the allowed values are continuos so are expressed as the min and max value
*/
type DefaultParser struct {
	// Allowed values for minutes
	minsValues []int
	// Allowed values for hours
	hoursValues []int
	// Allowed values for days of the month
	daysOfMonthValues []int
	// Allowed values for months
	monthsInt []int
	//Allowed values for days of the week
	daysOfWeekInt []int
	// A reference to the cron expression holder
	holder expressions.Holder
	// A field to keep the CronResults
	results *CronResults
	// A field to keed the CronElements
	cronElements *expressions.CronElements
}

// NewDefaultParser returns an instance of a default parser
func NewDefaultParser(expHolder expressions.Holder) (Parser, error) {
	dp := &DefaultParser{
		minsValues:        []int{0, 59},
		hoursValues:       []int{0, 23},
		daysOfMonthValues: []int{1, 31},
		monthsInt:         []int{1, 12},
		daysOfWeekInt:     []int{0, 6},
		holder:            expHolder,
	}
	ce, err := dp.holder.Elements()
	if err != nil {
		return nil, err
	}
	dp.cronElements = ce
	return dp, nil
}

/* parse method is the engine where the parsing logic is defined.
The method receives an input string (typically a field of a cron expression) and
an array of allowed values.
Every cron field can have multiple values that are separated by ','.
   After spliting by ',' each value can be one of the following
   - a single integer
   - an interval (format int-int)
   - * for all allowed values
   - / for a step expression
  The method call the right function to deal with the specific cases.
  The method return the list of value or an error if something goes wrong
*/
func (dp *DefaultParser) parse(input string, allowedValues []int) ([]int, error) {
	values := strings.Split(input, ",")
	results := []int{}
	for _, v := range values {
		// If we have an asterix we return all the allowed numbers.
		// We don't need to do any other calculation, for example an expression
		// *,1-4  because of the * we will select all the allowed values there is
		// no point to control the interval
		if string(v) == "*" {
			return utils.RangeValues(allowedValues), nil
		}

		// check if we have a step value
		steps := strings.Split(string(v), "/")
		if len(steps) == 2 {
			values, err := dp.manageSteps(steps, allowedValues)
			if err != nil {
				return results, err
			}
			results = append(results, values...)
			continue

		}
		// check if it is an interval
		interval := strings.Split(string(v), "-")
		if len(interval) == 2 {
			in, err := dp.manageIntervals(interval)
			if err != nil {
				return nil, err
			}
			values := utils.RangeValues(in)
			results = append(results, values...)
			continue
		}
		// It is a single value
		value, err := strconv.Atoi(string(v))
		if err != nil {
			return nil, err
		}
		results = append(results, value)
	}
	return results, nil
}

// We check that the two values passed are numbers and that the start value is less
// or equal of the end value.
func (dp *DefaultParser) manageIntervals(interval []string) ([]int, error) {
	results := []int{}
	start, err := strconv.Atoi(interval[0])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Syntax error for the interval %s", interval))
	}
	end, err := strconv.Atoi(interval[1])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Syntax error for the interval %s", interval))
	}
	if start > end {
		return nil, errors.New(fmt.Sprintf("Invalid interval, start is greater than end"))
	}
	results = append(results, start, end)
	return results, nil
}

/*
manageSteps is for cron fiels that are in the form of number/stepValue.
'number' can be an 'interval' or an '*'.
The input parameter is an array of 2 string that are the results of the the split with "/"
The allowedValues parameter is the list of allowed values for the field.
For the default parser the list of allowed values is made by a 2 int array where the first
element is the min value and the second one is the max value
*/
func (dp *DefaultParser) manageSteps(input []string, allowedValues []int) ([]int, error) {
	//The function finds the start value and create a series adding the step value until
	//we have reache the end value
	results := []int{}
	start := allowedValues[0]
	end := allowedValues[1]
	step, errStep := strconv.Atoi(input[1])
	if errStep != nil {
		return results, errors.New(fmt.Sprintf("Syntax error for step value"))
	}
	var err error

	// if we have an interval we take the start and end otherwise we have only a start
	interval := strings.Split(string(input[0]), "-")
	if len(interval) == 2 {
		in, errInt := dp.manageIntervals(interval)
		if errInt != nil {
			return results, errInt
		}
		start = in[0]
		end = in[1]

	} else if string(input[0]) != "*" {
		// We don't have an asterix and we don't have an interval, we only have a single number that is the start of the step
		start, err = strconv.Atoi(interval[0])
		if err != nil {
			return results, errors.New(fmt.Sprintf("Syntax error in %s", interval))
		}
	}

	for i := start; i <= end; i = i + step {
		results = append(results, i)
	}
	return results, nil

}

// generateResults is the internal method to return the results if they are already available
// or it generate all the results for all the expected fields.
func (dp *DefaultParser) generateResults() {
	if dp.results == nil {
		dp.results = &CronResults{}
	}
	dp.Minutes()
	dp.Hours()
	dp.DaysOfTheMonth()
	dp.DaysOfTheWeek()
	dp.Months()
	dp.Command()
}

// Minutes return the list of values for minutes or an error
func (dp *DefaultParser) Minutes() ([]int, error) {
	if dp.results != nil && len(dp.results.Minute) > 0 {
		return dp.results.Minute, nil
	}

	if dp.results == nil {
		dp.results = &CronResults{}
	}
	mins := dp.cronElements.Minute
	m, err := dp.parse(mins, dp.minsValues)
	if err != nil {
		return nil, err
	}
	// The results are as an array of int without duplicates and in ascending order
	dp.results.Minute = utils.SortedUniqueInts(m)
	// check if the values are in the allowed values, note that the check method requires a sorted array
	if !dp.inAllowedValues(dp.results.Minute, dp.minsValues) {
		return nil, errors.New(fmt.Sprintf("Minute value is not in the allowed interval %v\n", dp.minsValues))
	}
	return m, nil
}

// Hours return the list of values for hours or an error
func (dp *DefaultParser) Hours() ([]int, error) {
	if dp.results != nil && len(dp.results.Hour) > 0 {
		return dp.results.Hour, nil
	}

	if dp.results == nil {
		dp.results = &CronResults{}
	}

	hs := dp.cronElements.Hour
	h, err := dp.parse(hs, dp.hoursValues)

	if err != nil {
		return nil, err
	}
	// The results are as an array of int without duplicates and in ascending order
	dp.results.Hour = utils.SortedUniqueInts(h)
	// check if the values are in the allowed values, note that the check method requires a sorted array
	if !dp.inAllowedValues(dp.results.Hour, dp.hoursValues) {
		return nil, errors.New(fmt.Sprintf("Hour value is not in the allowed interval %v\n", dp.hoursValues))
	}
	return h, nil
}

// DaysOfTheMonth return the list of values for days of the month or an error
func (dp *DefaultParser) DaysOfTheMonth() ([]int, error) {
	if dp.results != nil && len(dp.results.DayMonth) > 0 {
		return dp.results.DayMonth, nil
	}

	if dp.results == nil {
		dp.results = &CronResults{}
	}

	dom := dp.cronElements.DayMonth
	dm, err := dp.parse(dom, dp.daysOfMonthValues)
	if err != nil {
		return nil, err
	}
	// The results are as an array of int without duplicates and in ascending order
	dp.results.DayMonth = utils.SortedUniqueInts(dm)
	// check if the values are in the allowed values, note that the check method requires a sorted array
	if !dp.inAllowedValues(dp.results.DayMonth, dp.daysOfMonthValues) {
		return nil, errors.New(fmt.Sprintf("Day of the Month value is not in the allowed interval %v\n", dp.daysOfMonthValues))
	}
	return dm, nil

}

// Months return the list of values for months or an error
func (dp *DefaultParser) Months() ([]int, error) {
	if dp.results != nil && len(dp.results.Month) > 0 {
		return dp.results.Month, nil
	}

	if dp.results == nil {
		dp.results = &CronResults{}
	}

	ms := dp.cronElements.Month
	m, err := dp.parse(ms, dp.monthsInt)
	if err != nil {
		return nil, err
	}
	// The results are as an array of int without duplicates and in ascending order
	dp.results.Month = utils.SortedUniqueInts(m)
	// check if the values are in the allowed values, note that the check method requires a sorted array
	if !dp.inAllowedValues(dp.results.Month, dp.monthsInt) {
		return nil, errors.New(fmt.Sprintf("Month value is not in the allowed interval %v\n", dp.monthsInt))
	}

	return m, nil

}

// DaysOfTheWeek return the list of values for days of the week or an error
func (dp *DefaultParser) DaysOfTheWeek() ([]int, error) {
	if dp.results != nil && len(dp.results.DayWeek) > 0 {
		return dp.results.DayWeek, nil
	}

	if dp.results == nil {
		dp.results = &CronResults{}
	}

	dow := dp.cronElements.DayWeek
	dw, err := dp.parse(dow, dp.daysOfWeekInt)
	if err != nil {
		return nil, err
	}
	// The results are as an array of int without duplicates and in ascending order
	dp.results.DayWeek = utils.SortedUniqueInts(dw)

	// check if the values are in the allowed values, note that the check method requires a sorted array
	if !dp.inAllowedValues(dp.results.DayWeek, dp.daysOfWeekInt) {
		return nil, errors.New(fmt.Sprintf("Day of the week value is not in the allowed interval %v\n", dp.daysOfWeekInt))
	}
	return dw, nil
}

//Command returns the command field or an error
func (dp *DefaultParser) Command() (string, error) {
	if dp.results != nil && len(dp.results.Command) > 0 {
		return dp.results.Command, nil
	}

	if dp.results == nil {
		dp.results = &CronResults{}
	}

	dp.results.Command = dp.cronElements.Command
	return dp.cronElements.Command, nil
}

// Results return the CronResults struct that contains the results of the parsing of a valid
// cron expression
func (dp *DefaultParser) Results() *CronResults {
	if dp.results == nil {
		dp.generateResults()
	}
	return dp.results
}

/*
Check if the values in an array are included in an interval.
the input array must be sorted in ascending order
the allowed values is an array of 2 values, where the first value is
the min value of the accepted values and the second one represents the max value
*/
func (dp *DefaultParser) inAllowedValues(input []int, allowedValues []int) bool {
	if input[0] < allowedValues[0] {
		return false
	}
	return input[len(input)-1] <= allowedValues[1]
}
