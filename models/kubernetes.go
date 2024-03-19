package models

type ElasticityPolicy struct {
	Type         string `json:"type"`
	Adjust       int    `json:"adjust"`
	Parameter    string `json:"parameter"`
	Operator     string `json:"operator"`
	Value        int    `json:"value"`
	PeriodNumber int    `json:"period_number"`
	Period       int    `json:"period"`
	Cooldown     int    `json:"cooldown"`
}

type ElasticityDict struct {
	Worker ElasticityWorker `json:"worker"`
}

type ElasticityWorker struct {
	MinVms             int                `json:"min_vms"`
	Cardinality        int                `json:"cardinality"`
	MaxVms             int                `json:"max_vms"`
	ElasticityPolicies []ElasticityPolicy `json:"elasticity_policies"`
}

type NodePool struct {
	Name             string         `json:"name"`
	SlugName         string         `json:"slug_name"`
	SKUID            string         `json:"sku_id"`
	SpecsName        string         `json:"specs_name"`
	WorkerNode       int            `json:"worker_node,omitempty"`
	ElasticityDict   ElasticityDict `json:"elasticity_dict,omitempty"`
	PolicyType       string         `json:"policy_type,omitempty"` //I changed it
	CustomParamName  string         `json:"custom_param_name,omitempty"`
	CustomParamValue string         `json:"custom_param_value,omitempty"`
}

type KubernetesCreate struct {
	Name      string     `json:"name"`
	SlugName  string     `json:"slug_name"`
	Version   string     `json:"version"`
	VPCID     string     `json:"vpc_id"`
	SKUID     string     `json:"sku_id"`
	NodePools []NodePool `json:"node_pools"`
}
