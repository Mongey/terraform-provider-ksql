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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the stream<Plug>_",
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

	return nil
}

func streamRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ksql.Client)
	name := d.Get("name").(string)
	log.Printf("[ERROR] Searching for stream %s", name)
	streams, err := c.ListStreams()
	if err != nil {
		return err
	}
	for _, s := range streams {
		//d.Set("query")
		log.Printf("[INFO] Found %s: %v", s.Name, s)
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
