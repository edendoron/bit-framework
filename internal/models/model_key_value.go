/*
 * bitTestResultsExporter
 *
 * This protocol defines the API for **test-results-exporter** service in the **BIT** functionality.
 *
 * API version: 1.0.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package models

type KeyValue struct {
	// key
	Key string `json:"key" validate:"required"`
	// value
	Value string `json:"value" validate:"required"`
}
