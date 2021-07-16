package models

type StorageExtendedStatus struct {
	// The type of the concrete used storage
	StorageType string `json:"storageType"`
	// The concrete DB used as storage. Relevant only if storageType='data_base'
	DataBase string `json:"dataBase,omitempty"`
}
