package goe2e

// DBConfig represents database configuration for DBaaS clusters.
type DBConfig struct {
	User        string `json:"user"`
	Password    string `json:"password"`
	Name        string `json:"name"`
	DBaaSNumber int    `json:"dbaas_number"`
}

// Software represents database software information.
type Software struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Engine  string `json:"engine"`
}

// DBNode represents a database node.
type DBNode struct {
	NodeName         string         `json:"node_name"`
	InstanceID       int            `json:"instance_id"`
	ClusterID        int            `json:"cluster_id"`
	NodeID           int            `json:"node_id"`
	VMID             int            `json:"vm_id"`
	Port             string         `json:"port"`
	PublicIPAddress  string         `json:"public_ip_address"`
	PrivateIPAddress string         `json:"private_ip_address"`
	AllowedIPs       AllowedIPs     `json:"allowed_ip_address"`
	ZabbixHostID     *int           `json:"zabbix_host_id"`
	Database         DBCreds        `json:"database"`
	RAM              string         `json:"ram"`
	CPU              int            `json:"cpu"`
	Disk             string         `json:"disk"`
	Status           string         `json:"status"`
	DBStatus         string         `json:"db_status"`
	CreatedAt        string         `json:"created_at"`
	Plan             Plan           `json:"plan"`
	SSL              bool           `json:"ssl"`
	Domain           *string        `json:"domain"`
	PublicPort       *int           `json:"public_port"`
	CommittedInfo    any            `json:"committed_info"`
	CommittedDetails []CommittedSKU `json:"committed_details"`
}

// DBCreds represents database credentials.
type DBCreds struct {
	ID       int      `json:"id"`
	Username string   `json:"username"`
	Database string   `json:"database"`
	PGDetail PGDetail `json:"pg_detail"`
}

// AllowedIPs represents allowed IP addresses configuration.
type AllowedIPs struct {
	WhitelistedIPs      []string `json:"whitelisted_ips"`
	TempIPs             []string `json:"temp_ips"`
	WhitelistedIPsTags  []string `json:"whitelisted_ips_tags"`
	TempIPsTags         []string `json:"temp_ips_tags"`
	WhitelistingRunning bool     `json:"whitelisting_in_progress"`
}

// Plan represents a database plan.
type Plan struct {
	Name                   string         `json:"name"`
	Price                  string         `json:"price"`
	TemplateID             int            `json:"template_id"`
	RAM                    string         `json:"ram"`
	CPU                    string         `json:"cpu"`
	Disk                   string         `json:"disk"`
	Currency               string         `json:"currency"`
	Software               Software       `json:"software"`
	AvailableInventoryStat bool           `json:"available_inventory_status"`
	PricePerHour           float64        `json:"price_per_hour"`
	PricePerMonth          float64        `json:"price_per_month"`
	CommittedSKU           []CommittedSKU `json:"committed_sku"`
}

// CommittedSKU represents a committed SKU.
type CommittedSKU struct {
	ID       int     `json:"committed_sku_id"`
	Name     string  `json:"committed_sku_name"`
	Message  string  `json:"committed_node_message"`
	Price    float64 `json:"committed_sku_price"`
	UptoDate string  `json:"committed_upto_date"`
	Days     int     `json:"committed_days"`
}

// PGDetail represents parameter group detail.
type PGDetail struct {
	ID int `json:"pg_id"`
}

// VPCMetadata represents VPC metadata for attach/detach operations.
type VPCMetadata struct {
	NetworkID string `json:"network_id"`
	VPCName   string `json:"vpc_name"`
	IPv4CIDR  string `json:"ipv4_cidr"`
}
