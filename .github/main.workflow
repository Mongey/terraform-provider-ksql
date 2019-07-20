workflow "New workflow" {
  on = "push"
  resolves = ["GitHub Action for Docker"]
}

action "Go FMT" {
  uses = "docker://golang:1.12.7"
  runs = "go"
  args = "fmt ./..."
}
