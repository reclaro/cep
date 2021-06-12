package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/reclaro/cep/expressions"
	"github.com/reclaro/cep/parsers"
	"github.com/reclaro/cep/printers"
)

/*
This script parses a cron string and expands each field to show the times at which it will run
*/
func main() {
	flag.Parse()

	if len(flag.Args()) > 1 {
		fmt.Println("The program accept only a single parameter as input string")
		os.Exit(1)
	}

	cmd := flag.Args()[0]
	if strings.Contains(cmd, "\n") {
		fmt.Println("The input string needs to be on a single line")
		os.Exit(1)
	}

	// we instantiate the expression holder that is responsible for checking the correctness of the cron expression string
	expressionHolder, err := expressions.NewDefaultSyntax(cmd)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// We instantiate the parser that is responsbile for parsing the string and expands all the fields
	p, err := parsers.NewDefaultParser(expressionHolder)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// We instantiate the printer that prints out the results based on a specific format/template
	prt := printers.NewSimple()
	res, err := p.Results()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	prt.Print(res)
}
