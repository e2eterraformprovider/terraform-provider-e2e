package models


type CreateContainerRegistryRequest struct {
	ProjectName string `json:"project_name"` 
	PreventVul  string `json:"prevent_vul"`  // "true" or "false"
	Severity    string `json:"severity"`     // "low" | "medium" | "high" | "critical"
}

type CreateContainerRegistryResponse struct {
	Code    int                         `json:"code"`
	Data    CreateContainerRegistryData `json:"data"`
	Errors  map[string]interface{}      `json:"errors"`
	Message string                      `json:"message"`
}

type CreateContainerRegistryData struct {
	SetupStatus string `json:"setup_status"`
}

type GetContainerRegistryProjectsResponse struct {
	Code            int                        `json:"code"`
	Data            []ContainerRegistryProject `json:"data"`
	Errors          map[string]interface{}     `json:"errors"`
	Message         string                     `json:"message"`
	TotalPageNumber int                        `json:"total_page_number"`
	TotalCount      int                        `json:"total_count"`
}

type ContainerRegistryProject struct {
	ID               int     `json:"id"`
	ProjectSize      float64 `json:"project_size"`
	DomainName       string  `json:"domain_name"`
	PreventVul       bool    `json:"prevent_vul"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
	Deleted          bool    `json:"deleted"`
	ProjectName      string  `json:"project_name"`
	ProjectID        int     `json:"project_id"`
	IsPublic         bool    `json:"is_public"`
	Severity         string  `json:"severity"`
	DeletedAt        *string `json:"deleted_at"`
	StorageLimit     int     `json:"storage_limit"`
	Location         string  `json:"location"`
	State            string  `json:"state"`
	Customer         int     `json:"customer"`
	MyAccountProject int     `json:"my_account_project"`
}

type DeleteContainerRegistryResponse struct {
	Code    int                         `json:"code"`
	Data    DeleteContainerRegistryData `json:"data"`
	Errors  map[string]interface{}      `json:"errors"`
	Message string                      `json:"message"`
}

type DeleteContainerRegistryData struct {
	Status string `json:"status"`
}
