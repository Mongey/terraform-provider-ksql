package ksql

import (
	"fmt"
	"log"

	"github.com/Mongey/ksql/ksql"
	"github.com/hashicorp/terraform/helper/schema"
)

func ksqlStreamResource() *schema.Resource {
	return &schema.Resource{
		Create: streamCreate,
		Read:   streamRead,
		Delete: streamDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "The name of the stream<Plug>_",
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

func streamCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	query := d.Get("query").(string)
	log.Printf("[WARN] Creating a Stream: %s with %s", name, query)
	c := meta.(*ksql.Client)
	q := fmt.Sprintf("CREATE STREAM %s %s", name, query)
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

	return streamRead(d, meta)
}

func streamRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ksql.Client)
	name := d.Get("name").(string)
	log.Printf("[ERROR] Searching for stream %s", name)
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

func streamDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ksql.Client)
	name := d.Get("name").(string)
	log.Printf("[INFO] Deleting stream %s", name)
	err := c.DropStream(&ksql.DropStreamRequest{Name: name})
	return err
}
