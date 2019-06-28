package ksql

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
)

var StreamResource = &Resource{"STREAM"}
var TableResource = &Resource{"TABLE"}

type Resource struct {
	Type string
}

// DiffSuppressEquivalentQueries returns true if the two compared queries are equivalent.
func  (r *Resource) DiffSuppressEquivalentQueries(k, old, new string, d *schema.ResourceData) bool {
	name := d.Get("name").(string)

	serverValue := MinifyKSQL(old)
	localValue := MinifyKSQL(r.FormatCreateQuery(name, new))
	
	if serverValue == localValue {
		return true
	}
	return false
}

func (r *Resource) FormatCreateQuery(name, query string) string {
	return fmt.Sprintf("CREATE %s %s %s", r.Type, name, query)
}
