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
