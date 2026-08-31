package models

type PlanResponse struct {
	Code    int            `json:"code"`
	Data    PlanData       `json:"data"`
	Errors  map[string]any `json:"errors"`
	Message string         `json:"message"`
}

type PlanData struct {
	TemplatePlans   []PlanTemplate     `json:"template_plans"`
	DatabaseEngines []EngineDefinition `json:"database_engines"`
}

type PostgressPlanUpgradeAction struct {
	TemplateID int `json:"template_id"`
}

type PostgressDiskAction struct {
	Size int `json:"size"`
}

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

type TemplateSoftware struct {
	SoftwareName    string `json:"name"`
	SoftwareVersion string `json:"version"`
	SoftwareEngine  string `json:"engine"`
}

type PlanCommittedSKU struct {
	SKUID           int     `json:"committed_sku_id"`
	SKUName         string  `json:"committed_sku_name"`
	SKUNodeMessage  string  `json:"committed_node_message"`
	SKUPrice        float64 `json:"committed_sku_price"`
	SKUEndDate      string  `json:"committed_upto_date"`
	SKUDurationDays int     `json:"committed_days"`
}

type EngineDefinition struct {
	EngineID          int     `json:"id"`
	EngineName        string  `json:"name"`
	EngineVersion     string  `json:"version"`
	EngineType        string  `json:"engine"`
	EngineDescription *string `json:"description"`
}
