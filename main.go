package main

import (
	"github.com/Mongey/terraform-provider-ksql/ksql"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: ksql.Provider})
}
