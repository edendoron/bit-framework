package models

type ReportBody struct {
	// Multiple tests reports set
	Reports []TestReport `json:"reports" validate:"required"`
}
