provider "ksql" {
  url = "http://localhost:8088"
}

resource "ksql_table" "logins" {
  name = "suspicious_attempts"

  query = "SELECT ip, count(*) FROM vault_failed_login_attempts_flattened WINDOW TUMBLING (size 30 second) GROUP BY ip having count(*) > 5"
}

resource "ksql_stream" "wp" {
  name = "vault_failed_login_attempts"

  query = "SELECT * FROM vault_logs WHERE type = 'response' AND response->data->error != '' AND request->path LIKE 'auth/userpass%';"
}
