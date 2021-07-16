package models

type KeyValue struct {
	// key
	Key string `json:"key" validate:"required"`
	// value
	Value string `json:"value" validate:"required"`
}
