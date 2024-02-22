package models

type BlockStorageCreate struct {
	Name string  `json:"name"`
	Size float64 `json:"size"`
	IOPS int     `json:"iops"`
	// ProjectID int    `json:"project_id"`
	// Location  string `json:"location"`
}

type BlockStorageResponse struct {
	Code    int                    `json:"code"`
	Data    []BlockStorage         `json:"data"`
	Errors  map[string]interface{} `json:"errors"`
	Message string                 `json:"message"`
}

// type BlockStorage struct {
// 	ID           int         `json:"id"`
// 	ImageName    string      `json:"image_name"`
// 	ResourceType interface{} `json:"resource_type"`
// 	LabelID      interface{} `json:"label_id"`
// }

type BlockAction struct {
}

// Just Trying
// type BlockStorageParams struct {
// 	ProjectID int    `json:"project_id"`
// 	Location  string `json:"location"`
// }

type BlockStorage struct {
	BlockID  int    `json:"block_id"`
	Name     string `json:"name"`
	Size     int    `json:"size"`
	Status   string `json:"status"`
	Template struct {
		DevPrefix    string `json:"DEV_PREFIX"`
		Driver       string `json:"DRIVER"`
		TotalIOPSSec string `json:"TOTAL_IOPS_SEC"`
	} `json:"template"`
	VMDetail  map[string]interface{} `json:"vm_detail"`
	CreatedOn string                 `json:"created_on"`
}
