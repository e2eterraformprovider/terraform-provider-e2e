package models

type SfsCreate struct{
	Name                    string        `json:"efs_name"`
	Plan                    string        `json:"efs_plan_name"`
	Vpc_id                  string         `json:"vpc_id"`
	Disk_size               int           `json:"efs_disk_size"`
	Disk_iops               int           `json:"efs_disk_iops"`
}