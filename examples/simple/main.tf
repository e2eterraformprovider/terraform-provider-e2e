# Simple example: Create a basic compute node

resource "e2e_node" "example" {
  name       = "example-node"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"
  project_id = var.project_id
  location   = var.location
}

output "node_id" {
  description = "The ID of the created node"
  value       = e2e_node.example.vm_id
}

output "node_name" {
  description = "The name of the created node"
  value       = e2e_node.example.name
}

output "node_status" {
  description = "The status of the created node"
  value       = e2e_node.example.status
}

