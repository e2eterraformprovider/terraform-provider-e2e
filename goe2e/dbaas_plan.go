package goe2e

// PlanResponse represents the response structure for DBaaS plan queries
type PlanResponse struct {
	Code    int            `json:"code"`
	Data    PlanData       `json:"data"`
	Errors  map[string]any `json:"errors"`
	Message string         `json:"message"`
}

// PlanData represents the data structure within PlanResponse
type PlanData struct {
	TemplatePlans   []PlanTemplate     `json:"template_plans"`
	DatabaseEngines []EngineDefinition `json:"database_engines"`
}

// PlanTemplate represents a database plan template
type PlanTemplate struct {
	PlanName             string             `json:"name"`
	PlanDisplayPrice     string             `json:"price"`
	PlanTemplateID       int                `json:"template_id"`
	PlanRAMGB            string             `json:"ram"`
	PlanCPUCores         string             `json:"cpu"`
	PlanDiskGB           string             `json:"disk"`
	PlanCurrency         string             `json:"currency"`
	PlanSoftware         TemplateSoftware   `json:"software"`
	IsInventoryAvailable bool               `json:"available_inventory_status"`
	PlanHourlyPrice      float64            `json:"price_per_hour"`
	PlanMonthlyPrice     float64            `json:"price_per_month"`
	CommittedSKUs        []PlanCommittedSKU `json:"committed_sku"`
}

// TemplateSoftware represents software information in a plan template
type TemplateSoftware struct {
	SoftwareName    string `json:"name"`
	SoftwareVersion string `json:"version"`
	SoftwareEngine  string `json:"engine"`
}

// PlanCommittedSKU represents a committed SKU in a plan
type PlanCommittedSKU struct {
	SKUID           int     `json:"committed_sku_id"`
	SKUName         string  `json:"committed_sku_name"`
	SKUNodeMessage  string  `json:"committed_node_message"`
	SKUPrice        float64 `json:"committed_sku_price"`
	SKUEndDate      string  `json:"committed_upto_date"`
	SKUDurationDays int     `json:"committed_days"`
}

// EngineDefinition represents a database engine definition
type EngineDefinition struct {
	EngineID          int     `json:"id"`
	EngineName        string  `json:"name"`
	EngineVersion     string  `json:"version"`
	EngineType        string  `json:"engine"`
	EngineDescription *string `json:"description"`
}
