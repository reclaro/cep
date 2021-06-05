package printers

import (
	"fmt"
	"github.com/reclaro/cep/parsers"
	"os"
	"strings"
	"text/template"
)

const (
	minutes    = "minute"
	hour       = "hour"
	dayOfMonth = "day of month"
	month      = "month"
	dayOfWeek  = "day of week"
	command    = "command"
)

const (
	table = `
{{.Minutes}}
{{.Hours}}
{{.DayMonth}}
{{.Month}}
{{.DayWeek}}
{{.Command}}
`
)

type Printer interface {
	Print(*parsers.CronResults)
}

type Simple struct {
	Minutes  string
	Hours    string
	DayMonth string
	Month    string
	DayWeek  string
	Command  string
}

func NewSimple() Printer {
	return &Simple{}
}

func (p *Simple) Print(exp *parsers.CronResults) {

	t := template.Must(template.New("Table").Parse(table))
	p.Minutes = fmt.Sprintf("%-14s%s", p.trimCol(minutes), strings.Trim(fmt.Sprintf("%+v", exp.Minute), "[]"))
	p.Hours = fmt.Sprintf("%-14s%s", p.trimCol(hour), strings.Trim(fmt.Sprintf("%+v", exp.Hour), "[]"))
	p.Month = fmt.Sprintf("%-14s%s", p.trimCol(month), strings.Trim(fmt.Sprintf("%+v", exp.Month), "[]"))
	p.DayMonth = fmt.Sprintf("%-14s%s", p.trimCol(dayOfMonth), strings.Trim(fmt.Sprintf("%+v", exp.DayMonth), "[]"))
	p.DayWeek = fmt.Sprintf("%-14s%s", p.trimCol(dayOfWeek), strings.Trim(fmt.Sprintf("%+v", exp.DayWeek), "[]"))
	p.Command = fmt.Sprintf("%-14s%s", p.trimCol(command), exp.Command)

	err := t.Execute(os.Stdout, p)
	if err != nil {
		fmt.Println("executing template:", err)
	}
}

func (p *Simple) trimCol(s string) string {
	return fmt.Sprintf("%.14s", s)
}
