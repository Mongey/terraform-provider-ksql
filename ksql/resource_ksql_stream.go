package ksql

import (
	"fmt"
	"log"
	"strings"

	"github.com/Mongey/ksql/ksql"
	"github.com/hashicorp/terraform/helper/schema"
)

func ksqlStreamResource() *schema.Resource {
	return &schema.Resource{
		Create: streamCreate,
		Read:   streamRead,
		Delete: streamDelete,
		Exists: streamExists,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "The name of the stream<Plug>_",
				DiffSuppressFunc: DiffSuppressCaseSensitivity,
			},
			"query": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "The query",
				DiffSuppressFunc: StreamResource.DiffSuppressEquivalentQueries,
			},
		},
	}
}

func streamCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	query := d.Get("query").(string)
	log.Printf("[WARN] Creating a Stream: %s with %s", name, query)
	c := meta.(*ksql.Client)
	q := StreamResource.FormatCreateQuery(name, query)
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
	err = d.Set("name", info.Name)
	if err != nil {
		return err
	}
	if !isSameCaseInsensitiveString(info.Type, StreamResource.Type) {
		return fmt.Errorf("incompatible type '%s' when expected '%s'", info.Type, StreamResource.Type)
	}
	if len(info.WriteQueries) > 0 {
		err = d.Set("query", info.WriteQueries[0].QueryString)
		if err != nil {
			return err
		}
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

func streamExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	c := meta.(*ksql.Client)
	name := d.Get("name").(string)
	log.Printf("[INFO] Looking if %s '%s' exists", strings.ToLower(StreamResource.Type), name)

	ls, err := c.ListStreams()
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
