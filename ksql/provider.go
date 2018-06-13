package ksql

import (
	"log"

	"github.com/Mongey/ksql/ksql"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KSQL_SERVER_URL", "http://localhost:8088"),
			},
		},

		ConfigureFunc: providerConfigure,
		ResourcesMap: map[string]*schema.Resource{
			"ksql_table":  ksqlTableResource(),
			"ksql_stream": ksqlStreamResource(),
		},
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	n := d.Get("url").(string)
	log.Printf("[INFO] configuring Provider %s", n)
	c := ksql.NewClient(n)

	return c, nil
}
