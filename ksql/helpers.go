package ksql

import (
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func isSameCaseInsensitiveString(a, b string) bool {
	if strings.ToLower(a) == strings.ToLower(b) {
		return true
	}
	return false
}

// DiffSuppressCaseSensitivity returns true if the two compared strings are the same, ignoring case sensitivity.
func DiffSuppressCaseSensitivity(k, old, new string, d *schema.ResourceData) bool {
	if isSameCaseInsensitiveString(old, new) {
		return true
	}
	return false
}

func MinifyKSQL(query string) string {
	// TODO: Add more KSQL minifying statements
	// The idea is that 2 statements with different text formatting compared against their output
	// from this function should now be equal if they represent the same logic.
	// String removal should include:
	//  - All line return
	//  - Duplicate spaces except those within quotation marks
	//  - Spaces after the KSQL operators "(", ")", "<", ">", ",", ...
	//  - All the other possible same things

	query = strings.Trim(query, "\n 	") // \n\s\t
	query = strings.ReplaceAll(query, "\n", "")
	return query
}