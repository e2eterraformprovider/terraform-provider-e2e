package models

type ObjectStore struct {
	ID                           float64 `json:"id"`
	Name                         string  `json:"name"`
	Status                       string  `json:"status"`
	BucketSize                   string  `json:"bucket_size"`
	CreatedOn                    string  `json:"created_on"`
	VersioningStatus             string  `json:"versioning_status"`
	LifecycleConfigurationStatus bool    `json:"lifecycle_configuration_status"`
}

type ObjectStorePayload struct {
}
