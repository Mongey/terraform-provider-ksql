package ksql

import (
	"fmt"
	"log"
	"testing"

	"github.com/Mongey/terraform-provider-kafka/kafka"
	r "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestBasicStream(t *testing.T) {
	topic, err := createTopic("vault")
	defer func() {
		err := topic.Delete()
		if err != nil {
			log.Printf("[ERROR] Unable to delete topic '%s': %v", topic.name, err)
		}
	}()
	if err != nil {
		t.Fatalf("[DEBUG] Could not create topic '%s': %v", topic.name, err)
	}
	r.Test(t, r.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []r.TestStep{
			{
				Config: testKSQLStreamQuery,
				Check:  testResourceStream_Check,
			},
		},
	})
}

func testResourceStream_Check(s *terraform.State) error {
	resourceState := s.Modules[0].Resources["ksql_stream.example"]
	if resourceState == nil {
		return fmt.Errorf("resource not found in state")
	}

	instanceState := resourceState.Primary
	if instanceState == nil {
		return fmt.Errorf("resource has no primary instance")
	}

	name := instanceState.ID

	if isSameCaseInsensitiveString(name, instanceState.Attributes["name"]) {
		return fmt.Errorf("id '%s' doesn't match expected name '%s'", name, instanceState.Attributes["name"])
	}

	if isSameCaseInsensitiveString(name, "vault_logs") {
		return fmt.Errorf("unexpected stream name '%s'", name)
	}

	return nil
}

const testKSQLStreamQuery = `
provider "ksql" {
	url = "http://localhost:8088"
}

resource "ksql_stream" "example" {
	name = "vault_logs"
	query = "(time VARCHAR, type VARCHAR, auth STRUCT<client_token VARCHAR, accessor VARCHAR, display_name VARCHAR, policies ARRAY<STRING>, token_policies ARRAY<STRING>, entity_id VARCHAR, token_type VARCHAR>, request STRUCT<id VARCHAR, operation VARCHAR, path VARCHAR, remote_address VARCHAR>, response STRUCT<data STRUCT<error VARCHAR>>, error VARCHAR) WITH (KAFKA_TOPIC='vault', VALUE_FORMAT='JSON', TIMESTAMP='time', TIMESTAMP_FORMAT='yyyy-MM-dd''T''HH:mm:ss[.SSSSSS][.SSSSS][.SSSS][.SSS][.SS][.S]''Z''');"
}
`

type kafkaTopic struct {
	name string
	client *kafka.Client
}

func (t *kafkaTopic) Delete() error {
	return t.client.DeleteTopic(t.name)
}

func createTopic(name string) (*kafkaTopic, error) {
	kafkaConfig := &kafka.Config{
		BootstrapServers: &[]string{"localhost:9092"},
		Timeout:          900,
	}
	kAdmin, err := kafka.NewClient(kafkaConfig)
	kTopic := &kafkaTopic{name, kAdmin}
	if err == nil {
		topic := kafka.Topic{
			Name:              name,
			Partitions:        1,
			ReplicationFactor: 1,
		}
		err = kAdmin.CreateTopic(topic)

		if err != nil {
			log.Printf("[ERROR] Creating Topic: %v", err)
			return kTopic, err
		}
	} else {
		log.Printf("[ERROR] Unable to create client: %s", err)
	}
	return kTopic, err
}
