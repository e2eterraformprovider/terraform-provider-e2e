package models

type SecurityGroupsResponse struct {
	Code    int             `json:"code"`
	Data    []SecurityGroup `json:"data"`
	Error   []interface{}   `json:"error"`
	Message string          `json:"message"`
}
type SecurityGroup struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Is_default  bool   `json:"is_default"`
	Rules       []Rule `json:"rules"`
}

type SecurityGroupCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Rules       []Rule `json:"rules"`
	Default     bool   `json:"default"`
}
type SecurityGroupUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Rules       []Rule `json:"rules"`
}

type Rule struct {
	Id            *int   `json:"id,omitempty"`
	Rule_type     string `json:"rule_type"`
	Protocol_name string `json:"protocol_name"`
	Port_range    string `json:"port_range"`
	Network       string `json:"network"`
	Network_cidr  string `json:"network_cidr,omitempty"`
	Network_size  *int   `json:"network_size,omitempty"`
	Description   string `json:"description"`
}
