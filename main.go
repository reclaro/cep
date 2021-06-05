package main

import (
	"flag"
	"fmt"
	"github.com/reclaro/cep/validators"
	"os"
)

/*
   * any value
   , value list separtor
   - range of values
   / step values
   First is minute:  allowedValues 0-59
   Second is hour:  allowed values 0-23
   Third is day of the month: allowed values 1-31
   Fourth is month: allowed Values 1-12 JAN-DEC
   Fifth day of the week: allowed Values 0-6 SUN-SAT
*/

func main() {
	flag.Parse()
	// TODO Check that we have a single string
	cmd := flag.Args()[0]

	fmt.Println("Cron string:", cmd)
	// TODO the validator can return a struct with the different tokens and be a single
	// struct with the fields and the parsing logic
	v := validators.NewDefaultSyntax(cmd)
	err := v.ValidateInput()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
