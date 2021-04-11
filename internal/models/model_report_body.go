/*
 * bit-test-results-exporter
 *
 * This protocol defines the API for **test-results-exporter** service in the **BIT** functionality.
 *
 * API version: 1.0.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package models

type ReportBody struct {
	// Multiple tests reports set
	Reports []TestReport `json:"reports" validate:"required"`
}