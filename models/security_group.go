package models

type SecurityGroupsResponse struct {
	Code    int             `json:"code"`
	Data    []SecurityGroup `json:"data"`
	Error   []interface{}   `json:"error"`
	Message string          `json:"message"`
}
type SecurityGroup struct {
	Id          float64 `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Is_default  bool    `json:"is_default"`
	Rules       []Rule  `json:"rules"`
}

type Rule struct {
	Id             float64 `json:"id"`
	Deleted        bool    `json:"deleted"`
	Created_at     string  `json:"created_at"`
	Updated_at     string  `json:"updated_at"`
	Rule_type      string  `json:"rule_type"`
	Protocol_name  string  `json:"protocol_name"`
	Port_range     string  `json:"port_range"`
	Network        string  `json:"network"`
	Is_active      bool    `json:"is_active"`
	Network_cidr   string  `json:"network_cidr"`
	Network_size   float64 `json:"network_size"`
	Security_group float64 `json:"security_group"`
}
