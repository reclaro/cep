package expressions

import (
	"errors"
	"fmt"
	"github.com/reclaro/cep/utils"
	"regexp"
	"strings"
)

/* Holder defines 2 methods, one is to return the different elements of a cron command from a cron expressions.
   The second method is to validate a cron expression depending on the syntax logic implemented by the struct implementing
   the interface.
*/
type Holder interface {
	Elements() (*CronElements, error)
	ValidateExpression(string) error
}

// CronElements is the struct for the 6 different fields that are present in a cron expression
type CronElements struct {
	Minute   string
	Hour     string
	DayMonth string
	Month    string
	DayWeek  string
	Command  string
}

// Default represents a default expression holder for the default cron job syntax.
// The default uses 5 required parameters separated by a white space
type DefaultSyntax struct {
	name                 string
	fields               int
	separator            string
	cronElements         *CronElements
	input                string
	daysMapper           map[string]string
	monthsMapper         map[string]string
	tokenValidatorString string
	regExpTokenValidator *regexp.Regexp
}

const (
	// fields is the value for the expected fields
	fields = 6
	// separator is the separator used in the input string.
	separator = " "
)

/*
NewDefaultSyntas implements the Holder interface.
   It return a new cron expression holder or error.
   The DefaultSyntax accepts a string input where the fields are separated by a single space
   For Days of the week it is possible to pass integer in the interval 0-6 where 0 is Sunday and it accepts also the
   following values: SUN, MON, TUE, WED, THU, FRI, SAT
   For Months is possible to pass integer in the interval 1-12
   Each field can be one of the following:
   int | int-int | * | * /int | int/int| int-int/int
*/
func NewDefaultSyntax(input string) (Holder, error) {
	ds := &DefaultSyntax{
		name:       "Standard Cron Expression",
		fields:     fields,
		separator:  separator,
		input:      input,
		daysMapper: map[string]string{"SUN": "0", "MON": "1", "TUE": "2", "WED": "3", "THU": "4", "FRI": "5", "SAT": "6"},
		monthsMapper: map[string]string{"JAN": "1",
			"FEB": "2",
			"MAR": "3",
			"APR": "4",
			"MAY": "5",
			"JUN": "6",
			"JUL": "7",
			"AUG": "8",
			"SEP": "9",
			"OCT": "10",
			"NOV": "11",
			"DEC": "12"},
		// Each bit can be in one of the following formats
		// int | int-int | * | */int | int/int| int-int/int
		tokenValidatorString: `^[0-9]+$|^\*$|^[0-9]+\-[0-9]+$|\*\/[0-9]+$|^[0-9]+\/[0-9]+$|^[0-9]+\-[0-9]+\/[0-9]+$`,
	}

	ds.regExpTokenValidator = regexp.MustCompile(ds.tokenValidatorString)
	return ds, nil
}

//ValidateExpression receives an input string and return an error if the syntax is not correct
func (ds *DefaultSyntax) ValidateExpression(input string) error {
	err := ds.validateFields(input)

	if err != nil {
		return err
	}
	return ds.tokenize()
}

// validateFields check that input string is made by the specific number of fields separated by a
// specific separator
func (ds *DefaultSyntax) validateFields(input string) error {
	tokens := strings.Split(input, ds.separator)
	if len(tokens) != ds.fields {
		return errors.New(fmt.Sprintf("Number of fields incorrect for %s, found %d and expected %d", ds.name, len(tokens), ds.fields))
	}
	return nil
}

// Elements return the cron string separated by each field or error if the input string is invalid
func (ds *DefaultSyntax) Elements() (*CronElements, error) {
	if ds.cronElements != nil {
		return ds.cronElements, nil
	}
	err := ds.validateFields(ds.input)
	if err != nil {
		return nil, err
	}
	err = ds.tokenize()
	if err != nil {
		return nil, err
	}
	return ds.cronElements, nil
}

// tokenize split the cron expression in the different fields. For each field we validate the syntax,
// if it is not correct we return an error
func (ds *DefaultSyntax) tokenize() error {
	// Note this split on a single white space, if we want to split on white spaces
	// we can use a regexp for it regexp.MustCompile(`\S+`) and then re.FindAllString(input, -1)
	tokens := strings.Split(ds.input, ds.separator)
	// We check if it has been passed the strings format for Day of week and month
	// and we convert it to the integers
	tokens[3] = utils.StringToNumber(tokens[3], ds.monthsMapper)
	tokens[4] = utils.StringToNumber(tokens[4], ds.daysMapper)

	// we do not validate the command token that is in the last position
	for i := 0; i < len(tokens)-2; i++ {
		err := ds.validateTokens(tokens[i])
		if err != nil {
			return err
		}
	}
	//TO DO use "enum" instead of integer here
	ce := &CronElements{Minute: tokens[0],
		Hour:     tokens[1],
		DayMonth: tokens[2],
		Month:    tokens[3],
		DayWeek:  tokens[4],
		Command:  tokens[5],
	}
	ds.cronElements = ce
	return nil
}

// validteTokens receive a fields as a string and it validates the correct syntax. It retunrs an error
// if the syntax is not valid.
func (ds *DefaultSyntax) validateTokens(token string) error {
	// split each token on the comma
	t := strings.Split(token, ",")
	for _, str := range t {
		isValid := ds.regExpTokenValidator.MatchString(str)
		if !isValid {
			return errors.New(fmt.Sprintf("Invalid input string '%s' please check the correct syntax", ds.input))
		}
	}
	return nil
}
