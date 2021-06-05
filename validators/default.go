package validators

import (
	"errors"
	"fmt"
	"strings"
)

type Validator interface {
	ValidateInput() error
	CronExpression() (*CronElements, error)
}

type CronElements struct {
	Minute   string
	Hour     string
	DayMonth string
	Month    string
	DayWeek  string
	Command  string
}

// DefaultSyntax represents a default validator for the cron job syntax.
// The default uses 5 required parameters separated by a white space
type DefaultSyntax struct {
	name         string
	fields       int
	separator    string
	cronElements *CronElements
	input        string
}

const (
	fields    = 6
	separator = " "
)

func NewDefaultSyntax(input string) Validator {
	return &DefaultSyntax{
		name:      "Standard Cron Expression",
		fields:    fields,
		separator: separator,
		input:     input,
	}
}

func (ds *DefaultSyntax) ValidateInput() error {
	tokens := strings.Split(ds.input, ds.separator)
	if len(tokens) != ds.fields {
		return errors.New(fmt.Sprintf("Number of fields incorrect for %s, found %d and expected %d", ds.name, len(tokens), ds.fields))
	}
	return nil
}

func (ds *DefaultSyntax) CronExpression() (*CronElements, error) {
	if ds.cronElements != nil {
		return ds.cronElements, nil
	}
	err := ds.tokenize()
	if err != nil {
		return nil, err
	}
	return ds.cronElements, nil
}

func (ds *DefaultSyntax) tokenize() error {
	err := ds.ValidateInput()
	if err != nil {
		return err
	}
	tokens := strings.Split(ds.input, ds.separator)
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
