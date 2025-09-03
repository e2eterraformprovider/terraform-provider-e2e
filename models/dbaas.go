package models

type DBResponse struct {
	Code    int    `json:"code"`
	Data    DB     `json:"data"`
	Errors  any    `json:"errors"`
	Message string `json:"message"`
}

type DB struct {
	ID                  int      `json:"id"`
	Name                string   `json:"name"`
	Status              string   `json:"status"`
	StatusTitle         string   `json:"status_title"`
	StatusActions       []string `json:"status_actions"`
	NumInstances        int      `json:"num_instances"`
	Software            Software `json:"software"`
	MasterNode          DBNode   `json:"master_node"`
	ConnectivityDetail  string   `json:"connectivity_detail"`
	VectorDBStatus      string   `json:"vector_database_status"`
	ProjectName         string   `json:"project_name"`
	SnapshotExist       bool     `json:"snapshot_exist"`
	ZookeeperInstances  int      `json:"zookeeper_instances"`
	SlaveInstances      int      `json:"slave_instances"`
	IsEncryptionEnabled bool     `json:"isEncryptionEnabled"`
}

type Software struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Engine  string `json:"engine"`
}

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
	CPU              int         `json:"cpu"`
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

type DBCreds struct {
	ID       int      `json:"id"`
	Username string   `json:"username"`
	Database string   `json:"database"`
	PGDetail PGDetail `json:"pg_detail"`
}

type AllowedIPs struct {
	WhitelistedIPs      []string `json:"whitelisted_ips"`
	TempIPs             []string `json:"temp_ips"`
	WhitelistedIPsTags  []string `json:"whitelisted_ips_tags"`
	TempIPsTags         []string `json:"temp_ips_tags"`
	WhitelistingRunning bool     `json:"whitelisting_in_progress"`
}

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

type MySqlCreate struct {
	Name             string   `json:"name"`
	Database         DBConfig `json:"database"`
	Vpcs             []VPC    `json:"vpcs"`
	SoftwareID       int      `json:"software_id"`
	TemplateID       int      `json:"template_id"`
	ParameterGroupId int      `json:"pg_id,omitempty"`
	PublicIPRequired bool     `json:"public_ip_required"`
	Group            string   `json:"group"`
}

type CommittedSKU struct {
	ID       int     `json:"committed_sku_id"`
	Name     string  `json:"committed_sku_name"`
	Message  string  `json:"committed_node_message"`
	Price    float64 `json:"committed_sku_price"`
	UptoDate string  `json:"committed_upto_date"`
	Days     int     `json:"committed_days"`
}

type MariaDBCreateRequest struct {
	Name                 string        `json:"name"`
	SoftwareID           int           `json:"software_id"`
	TemplateID           int           `json:"template_id"`
	PublicIPRequired     bool          `json:"public_ip_required"`
	Group                string        `json:"group"`
	VPCs                 []VPCMetadata `json:"vpcs,omitempty"`
	Database             DBConfig      `json:"database"`
	PGID                 int           `json:"pg_id"`
	IsEncryptionEnabled  bool          `json:"isEncryptionEnabled"`
	EncryptionPassphrase string        `json:"encryption_passphrase"`
}

type DBCreateRequest struct {
	Name                string   `json:"name"`
	SoftwareID          int      `json:"software_id"`
	TemplateID          int      `json:"template_id"`
	PublicIPRequired    bool     `json:"public_ip_required"`
	Group               string   `json:"group"`
	VPCs                []VPC    `json:"vpcs"`
	Database            DBConfig `json:"database"`
	PGID                *int     `json:"pg_id,omitempty"`
	IsEncryptionEnabled bool     `json:"isEncryptionEnabled"`
}

type DBConfig struct {
	User        string `json:"user"`
	Password    string `json:"password"`
	Name        string `json:"name"`
	DBaaSNumber int    `json:"dbaas_number"`
}

type AttachDetachVPCRequest struct {
	Action string        `json:"action"`
	VPCs   []VPCMetadata `json:"vpcs"`
}

type VPCMetadata struct {
	NetworkID string `json:"network_id"`
	VPCName   string `json:"vpc_name"`
	IPv4CIDR  string `json:"ipv4_cidr"`
}

type ParameterGroupRequest struct {
	Action string `json:"action"`
}

type UpgradePlanRequest struct {
	TemplateID int `json:"template_id"`
}

type DiskUpgradeRequest struct {
	Size int `json:"size"`
}

type PGDetail struct {
	ID int `json:"pg_id"`
}

type MySQlPlanUpgradeAction struct {
	TemplateID int `json:"template_id"`
}

type MYSQLExpandDisk struct {
	Size int `json:"size"`
}

type AttachVPCPayloadRequest struct {
	Action string `json:"action"`
	VPCs   []VPC  `json:"vpcs"`
}

type VPC struct {
	VpcName    string  `json:"vpc_name,omitempty"`
	Ipv4_cidr  string  `json:"ipv4_cidr,omitempty"`
	Network_id float64 `json:"network_id,omitempty"`
}
