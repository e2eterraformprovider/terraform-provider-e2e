# Simple Example

This example demonstrates how to create a basic compute node using the E2E Networks Terraform provider.

## Prerequisites

- Terraform installed (0.13+)
- E2E Networks account with API credentials
- Provider installed locally (see main README for installation instructions)

## Setup

1. Navigate to this example directory:

   ```bash
   cd examples/simple
   ```

2. Create a `terraform.tfvars` file with your credentials:

   ```hcl
   api_key    = "your-api-key"
   auth_token = "your-auth-token"
   project_id = "your-project-id"
   location   = "in-mumbai-1"
   ```

   **Note**: Alternatively, you can set these as environment variables:

   ```bash
   export TF_VAR_api_key="your-api-key"
   export TF_VAR_auth_token="your-auth-token"
   export TF_VAR_project_id="your-project-id"
   export TF_VAR_location="in-mumbai-1"
   ```

## Usage

1. Initialize Terraform:

   ```bash
   terraform init
   ```

2. Validate the configuration:

   ```bash
   terraform validate
   ```

3. Preview the changes (dry-run):

   ```bash
   terraform plan
   ```

4. Apply the configuration (creates real resources):

   ```bash
   terraform apply
   ```

5. Destroy the resources when done:
   ```bash
   terraform destroy
   ```

## What This Example Creates

- A compute node named "example-node"
- Plan: c2-2c-4gb (2 CPU cores, 4GB RAM)
- Image: Ubuntu 20.04
- Location: Mumbai (or as specified in variables)

## Important Notes

- This example creates real resources that may incur costs
- Always run `terraform destroy` when you're done testing
- Replace placeholder values with your actual E2E Networks credentials
- Make sure you have sufficient quota in your project
