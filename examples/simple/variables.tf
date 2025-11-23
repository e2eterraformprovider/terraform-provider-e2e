variable "api_key" {
  description = "E2E Networks API Key"
  type        = string
  sensitive   = true
}

variable "auth_token" {
  description = "E2E Networks Auth Token"
  type        = string
  sensitive   = true
}

variable "project_id" {
  description = "E2E Networks Project ID"
  type        = string
}

variable "location" {
  description = "E2E Networks Location (e.g., in-mumbai-1)"
  type        = string
  default     = "in-mumbai-1"
}

