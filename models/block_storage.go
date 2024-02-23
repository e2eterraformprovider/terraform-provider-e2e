package models

type BlockStorageCreate struct {
	Name string  `json:"name"`
	Size float64 `json:"size"`
	IOPS int     `json:"iops"`
}

type BlockStorageResponse struct {
	Code    int                    `json:"code"`
	Data    []BlockStorage         `json:"data"`
	Errors  map[string]interface{} `json:"errors"`
	Message string                 `json:"message"`
}

type BlockAction struct {
}


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