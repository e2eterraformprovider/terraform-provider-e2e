package models

// FaasNamespaceCreate represents the request to create a FaaS namespace
type FaasNamespaceCreate struct {
	Name string `json:"name"`
}

// FaasNamespaceResponse represents the response from namespace operations
type FaasNamespaceResponse struct {
	Code    int               `json:"code"`
	Data    FaasNamespaceData `json:"data"`
	Error   []interface{}     `json:"error"`
	Message string            `json:"message"`
}

type FaasNamespaceData struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

// FaasFunctionCreate represents the request to create a FaaS function
type FaasFunctionCreate struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Runtime     string            `json:"runtime"`
	Code        string            `json:"code"`
	MemoryMB    int               `json:"memory_mb,omitempty"`
	Timeout     int               `json:"timeout_seconds,omitempty"`
	MinReplicas int               `json:"min_replicas,omitempty"`
	MaxReplicas int               `json:"max_replicas,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
}

// FaasFunctionUpdate represents the request to update a FaaS function
type FaasFunctionUpdate struct {
	Code        string            `json:"code,omitempty"`
	MemoryMB    int               `json:"memory_mb,omitempty"`
	Timeout     int               `json:"timeout_seconds,omitempty"`
	MinReplicas int               `json:"min_replicas,omitempty"`
	MaxReplicas int               `json:"max_replicas,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
}

// FaasFunctionResponse represents the response from function operations
type FaasFunctionResponse struct {
	Code    int           `json:"code"`
	Data    FaasFunction  `json:"data"`
	Error   []interface{} `json:"error"`
	Message string        `json:"message"`
}

// FaasFunction represents a FaaS function
type FaasFunction struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Runtime     string            `json:"runtime"`
	Code        string            `json:"code,omitempty"`
	MemoryMB    int               `json:"memory_mb"`
	Timeout     int               `json:"timeout_seconds"`
	MinReplicas int               `json:"min_replicas"`
	MaxReplicas int               `json:"max_replicas"`
	Environment map[string]string `json:"environment"`
	EndpointURL string            `json:"endpoint_url"`
	Status      string            `json:"status"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
}

// FaasFunctionListResponse represents the response from listing functions
type FaasFunctionListResponse struct {
	Code    int            `json:"code"`
	Data    []FaasFunction `json:"data"`
	Error   []interface{}  `json:"error"`
	Message string         `json:"message"`
}

// FaasLogsResponse represents the response from function logs
type FaasLogsResponse struct {
	Code    int           `json:"code"`
	Data    []FaasLog     `json:"data"`
	Error   []interface{} `json:"error"`
	Message string        `json:"message"`
}

type FaasLog struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
	Level     string `json:"level"`
}
