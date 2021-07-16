package models

type Bandwidth struct {
	// Upper limit value, 0 or negative means unlimited
	Size float32 `json:"size,omitempty" validate:"required"`
	// KiB/MiB/GiB/TiB/K/M/G/T, 1_KiB = 1024, 1_K = 1000
	UnitsPerSecond string `json:"unitsPerSecond,omitempty" validate:"required"`
}
