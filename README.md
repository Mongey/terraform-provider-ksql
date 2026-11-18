# `terraform-provider-ksql`
[![CircleCI](https://circleci.com/gh/Mongey/terraform-provider-ksql.svg?style=svg&circle-token=320e9b975067221dd59cc169e83b8faf53ea5062)](https://circleci.com/gh/Mongey/terraform-provider-ksql)

A [Terraform][1] plugin for managing [Confluent KSQL Server][2].

## Contents

* [Installation](#installation)
  * [Developing](#developing)
* [`ksql` Provider](#provider-configuration)
* [Resources](#resources)
  * [`ksql_stream`](#ksql_stream)
  * [`ksql_table`](#ksql_table)

## Installation

Download and extract the [latest
release](https://github.com/Mongey/terraform-provider-ksql/releases/latest) to
your [terraform plugin directory][third-party-plugins] (typically `~/.terraform.d/plugins/`)

### Developing

Install [Go 1.27.4 or later][install-go] and clone the repository into any directory.

```sh
git clone https://github.com/Mongey/terraform-provider-ksql.git
cd terraform-provider-ksql
make build
make vet
make test
```

Dependencies are tracked in `vendor/`, which Go uses automatically. `make build`
writes the provider binary into the repository root. Ordinary tests skip Kafka/KSQL
acceptance tests. To run acceptance tests, start Kafka on `localhost:9092` and KSQL
Server on `localhost:8088`, then run `make testacc`. These tests create Kafka topics.

After changing dependencies, run `go mod tidy` and `go mod vendor`, then commit
`go.mod`, `go.sum` and `vendor/` together.

## Provider Configuration

### Example

```hcl
provider "ksql" {
  url = "http://localhost:8083"
}
```

## Resources
### `ksql_stream`

A resource for managing KSQL streams
```hcl
resource "ksql_stream" "actions" {
  name = "vip_actions"
  query = "SELECT userid, page, action
              FROM clickstream c
              LEFT JOIN users u ON c.userid = u.user_id
              WHERE u.level =
              'Platinum';"
}
```

### `ksql_table`

A resource for managing KSQL tables
```hcl
resource "ksql_table" "users" {
  name = "users-thing"
  query = "SELECT error_code,
            count(*),
            FROM monitoring_stream
            WINDOW TUMBLING (SIZE 1 MINUTE)
            WHERE  type = 'ERROR'
            GROUP BY error_code;"
  }
}
```


[install-go]: https://go.dev/doc/install
[1]: https://www.terraform.io
[2]: https://www.confluent.io/product/ksql/
