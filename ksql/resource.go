package ksql

var StreamResource = &Resource{"STREAM"}
var TableResource = &Resource{"TABLE"}

type Resource struct {
	Type string
}
