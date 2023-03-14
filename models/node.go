package models

type Node struct {
	Name              string `json:"name"`
	Label             string `json:"label"`
	Plan              string `json:"plan"`
	Backup            bool   `json:"backup"`
	Image             string `json:"image"`
	Default_public_ip bool   `json:"default_public_id"`
	Disable_password  bool   `json:"disable_password"`
	Enable_bitninja   bool   `json:"enable_bitninja"`
	Is_ipv6_availed   bool   `json:"is_ipv6_availed"`
	Is_saved_image    bool   `json:"is_saved_image"`
	Region            string `json:"region"`
	Reserve_ip        string `json:"reserve_ip"`
	Vpc_id            string `json:"vpc_id"`
	//Ngc_container_id        int      `json:"ngc_container_id"`
	//Saved_image_template_id int      `json:"saved_image_template_id"`
	Security_group_id int           `json:"security_group_id"`
	SSH_keys          []interface{} `json:"ssh_keys"`
}
type NodeAction struct {
	Type string `json:"type"`
}
