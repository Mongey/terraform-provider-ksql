package ksql

import (
	"fmt"
	"log"
	"testing"

	r "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestBasicTable(t *testing.T) {
	err := createTopic("users")
	if err != nil {
		log.Printf("[DEBUG] state %v", err)
		t.Fatalf("Could not create the topic: %s", err)
	}
	r.Test(t, r.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []r.TestStep{
			{
				Config: testKSQLTableQuery,
				Check:  testResourceTable_Check,
			},
		},
	})
}

func testResourceTable_Check(s *terraform.State) error {
	log.Printf("[ERROR] %v", s)
	resourceState := s.Modules[0].Resources["ksql_table.example"]
	if resourceState == nil {
		return fmt.Errorf("resource not found in state")
	}

	instanceState := resourceState.Primary
	if instanceState == nil {
		return fmt.Errorf("resource has no primary instance")
	}

	name := instanceState.ID

	if name != instanceState.Attributes["name"] {
		return fmt.Errorf("id doesn't match name")
	}

	if name != "users" {
		return fmt.Errorf("unexpected topic name %s", name)
	}

	//_ = testProvider.Meta().(*ksql.Client)

	return nil
}

const testKSQLTableQuery = `
provider "ksql" {
	url = "http://localhost:8088"
}

resource "ksql_table" "example" {
	name = "users"
	query = "(registertime BIGINT, userid VARCHAR, gender VARCHAR, regionid VARCHAR) WITH (KAFKA_TOPIC = 'users', VALUE_FORMAT='JSON', KEY = 'userid');"
}
`
