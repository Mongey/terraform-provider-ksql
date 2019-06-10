package ksql

import (
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

// DiffSuppressCaseSensitivity returns true if the two compared strings are the same, ignoring case sensitivity.
func DiffSuppressCaseSensitivity(k, old, new string, d *schema.ResourceData) bool {
	if strings.ToLower(old) == strings.ToLower(new) {
		return true
	}
	return false
}
