package ksql

import (
	"fmt"
	"log"

	"github.com/Mongey/ksql/ksql"
	"github.com/hashicorp/terraform/helper/schema"
)

func ksqlTableResource() *schema.Resource {
	return &schema.Resource{
		Create: tableCreate,
		Read:   tableRead,
		Delete: tableDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the table",
			},
			"query": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The query",
			},
		},
	}
}

func tableCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	query := d.Get("query").(string)
	log.Printf("[WARN] Creating a table: %s with %s", name, query)
	c := meta.(*ksql.Client)
	q := fmt.Sprintf("CREATE TABLE %s %s", name, query)
	log.Printf("[WARN] Query %s", q)

	r := ksql.Request{
		KSQL: q,
	}
	resp, err := c.Do(r)
	log.Printf("[RESP] %v", resp)
	if err != nil {
		return err
	}
	d.SetId(name)
	return nil
}

func tableRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ksql.Client)
	name := d.Get("name").(string)
	log.Printf("[ERROR] Searching for table %s", name)
	tables, err := c.ListTables()
	if err != nil {
		return err
	}
	for _, t := range tables {
		//d.Set("query")
		log.Printf("[INFO] Found %s: %v", t.Name, t)
	}
	return nil
}

func tableDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ksql.Client)
	name := d.Get("name").(string)
	log.Printf("[INFO] Deleting table %s", name)
	err := c.DropTable(&ksql.DropTableRequest{Name: name})
	return err
}
