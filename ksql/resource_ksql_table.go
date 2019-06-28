package ksql

import (
	"fmt"
	"log"
	"strings"

	"github.com/Mongey/ksql/ksql"
	"github.com/hashicorp/terraform/helper/schema"
)

func ksqlTableResource() *schema.Resource {
	return &schema.Resource{
		Create: tableCreate,
		Read:   tableRead,
		Delete: tableDelete,
		Exists: tableExists,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "The name of the table",
				DiffSuppressFunc: DiffSuppressCaseSensitivity,
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
	return tableRead(d, meta)
}

func tableRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ksql.Client)
	name := d.Get("name").(string)
	log.Printf("[ERROR] Searching for table %s", name)
	info, err := c.Describe(name)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Found %s: %v", info.Name, info)
	d.Set("name", info.Name)
	if len(info.WriteQueries) > 0 {
		d.Set("query", info.WriteQueries[0].QueryString)
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

func tableExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	c := meta.(*ksql.Client)
	name := d.Get("name").(string)
	log.Printf("[INFO] Looking if %s '%s' exists", strings.ToLower(TableResource.Type), name)

	ls, err := c.ListTables()
	if err != nil {
		return true, err
	}

	for _, r := range ls {
		log.Printf("[INFO] Found %s: %+v", r.Name, r)
		if isSameCaseInsensitiveString(r.Name, name) {
			return true, err
		}
	}

	return false, err
}
