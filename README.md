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

0. [Install go][install-go]
0. Clone repository to: `$GOPATH/src/github.com/Mongey/terraform-provider-ksql`
    ``` bash
    mkdir -p $GOPATH/src/github.com/Mongey/terraform-provider-ksql; cd $GOPATH/src/github.com/Mongey/
    git clone https://github.com/Mongey/terraform-provider-ksql.git
    ```
0. Build the provider `make build`
0. Run the tests `make test`

0. Build the provider `make build`

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


[install-go]: https://golang.org/doc/install#install
[1]: https://www.terraform.io
[2]: https://www.confluent.io/product/ksql/
