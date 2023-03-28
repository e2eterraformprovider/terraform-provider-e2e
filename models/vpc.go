package models

type VpcsResponse struct {
	Code    int           `json:"code"`
	Data    []Vpc         `json:"data"`
	Error   []interface{} `json:"error"`
	Message string        `json:"message"`
}
type Vpc struct {
	Created_at string  `json:"created_at"`
	State      string  `json:"state"`
	Name       string  `json:"name"`
	Ipv4_cidr  string  `json:"ipv4_cidr"`
	Network_id float64 `json:"network_id"`
	Gateway_ip string  `json:"gateway_ip"`
	Pool_size  float64 `json:"pool_size"`
	Is_active  bool    `json:"is_active"`
}
