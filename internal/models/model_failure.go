package models

type Failure struct {
	//
	Description FailureDescription `json:"unit-name" validate:"required"`
	//
	ExaminationRule FailureExaminationRule `json:"examination-rule" validate:"required"`
	//
	ReportDuration FailureReportDuration `json:"report-duration" validate:"required"`
	//
	Dependencies FailureDependencies `json:"dependencies" validate:"required"`
}

type FailureDependencies struct {
	//
	BelongsToGroup []string `json:"belongs-to-group" validate:"required"`
	//
	MasksOtherGroup []string `json:"masks-other-group" validate:"required"`
}
