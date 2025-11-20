<div align="center">
  <img src="https://e2enetworks.com/OnlyE2E.svg" alt="E2E Networks Logo" width="180"/>

  # 🚀 Terraform Provider for E2E Networks

  <p>
    <a href="https://cloud.e2enetworks.com/"><img src="https://img.shields.io/badge/E2E_Cloud-Console-blue?style=for-the-badge" alt="E2E Cloud Console"/></a>
    <a href="https://docs.e2enetworks.com/"><img src="https://img.shields.io/badge/E2E-Documentation-green?style=for-the-badge" alt="E2E Documentation"/></a>
    <a href="https://registry.terraform.io/providers/e2eterraformprovider/e2e/latest"><img src="https://img.shields.io/badge/Terraform-Provider-purple?style=for-the-badge" alt="Terraform E2E Provider"/></a>
    <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-yellow?style=for-the-badge" alt="License: MIT"/></a>
  </p>

  **A Production-Ready Infrastructure as Code Solution for E2E Networks Cloud**
</div>

---

## Overview

The E2E Networks Terraform Provider enables you to manage your E2E Networks cloud infrastructure using Terraform's infrastructure-as-code approach. Build, modify, and version your cloud infrastructure safely and efficiently.

### Supported Resources

- **Compute**: Virtual Machines (Nodes), Images, SSH Keys
- **Networking**: VPCs, Load Balancers, Reserved IPs, Security Groups
- **Storage**: Block Storage, Shared File Storage (SFS), Object Storage
- **Databases**: PostgreSQL, MySQL, MariaDB (DBaaS)
- **Containers**: Kubernetes, Container Registry
- **Serverless**: FaaS Functions ⚡ **NEW!**
- **Auto Scaling**: Scaler Groups

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.13+
- [Go](https://golang.org/doc/install) 1.18+ (for building the provider plugin)
- E2E Networks account with API credentials

## Using the Provider

### Installation

The provider is available on the [Terraform Registry](https://registry.terraform.io/providers/e2eterraformprovider/e2e/latest). Add it to your Terraform configuration:

```hcl
terraform {
  required_providers {
    e2e = {
      source  = "e2eterraformprovider/e2e"
      version = "~> 1.0"
    }
  }
}

provider "e2e" {
  api_key      = var.e2e_api_key
  auth_token   = var.e2e_auth_token
  api_endpoint = "https://api.e2enetworks.com/myaccount/api/v1/"
}
```

### Authentication

Configure your E2E Networks credentials using one of these methods:

**Option 1: Environment Variables**
```bash
export SERVICE_API_KEY="your-api-key"
export SERVICE_AUTH_TOKEN="your-auth-token"
```

**Option 2: Provider Configuration**
```hcl
provider "e2e" {
  api_key    = "your-api-key"
  auth_token = "your-auth-token"
}
```

### Quick Start Example

```hcl
# Create a VPC
resource "e2e_vpc" "main" {
  vpc_name   = "production-vpc"
  location   = "Delhi"
  project_id = var.project_id
}

# Deploy a virtual machine
resource "e2e_node" "web" {
  node_name       = "web-server-01"
  node_vcpu_count = "4"
  node_ram_mb     = "8192"
  node_image_name = "Ubuntu 22.04 LTS"
  plan_name       = "4C-8GB-80GB"
  location        = "Delhi"
  project_id      = var.project_id
  vpc_id          = e2e_vpc.main.id
}

# Create a FaaS function
resource "e2e_faas_function" "api" {
  name            = "order-processor"
  namespace       = "production"
  runtime         = "python-3.11-fastapi"
  code_inline     = file("${path.module}/handler.py")
  memory_mb       = 512
  timeout_seconds = 60
  min_replicas    = 1
  max_replicas    = 5

  environment = {
    DB_HOST = "postgres.example.com"
  }

  project_id = var.project_id
  location   = "Delhi"
}

output "function_url" {
  value = e2e_faas_function.api.endpoint_url
}
```

## Building The Provider

If you want to build the provider from source:

### Clone the Repository

```bash
mkdir -p $GOPATH/src/github.com/e2eterraformprovider
cd $GOPATH/src/github.com/e2eterraformprovider
git clone https://github.com/e2eterraformprovider/terraform-provider-e2e
cd terraform-provider-e2e
```

### Build

```bash
go build -o terraform-provider-e2e
```

### Install Locally

```bash
# Create local plugins directory
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/e2eterraformprovider/e2e/1.0.0/linux_amd64

# Copy the built provider
cp terraform-provider-e2e ~/.terraform.d/plugins/registry.terraform.io/e2eterraformprovider/e2e/1.0.0/linux_amd64/
```

## Documentation

For detailed documentation on all resources and data sources, visit:

- **[Official Provider Documentation](https://registry.terraform.io/providers/e2eterraformprovider/e2e/latest/docs)**
- **[E2E Networks API Documentation](https://docs.e2enetworks.com/api/myaccount/)**
- **[E2E Networks Cloud Console](https://cloud.e2enetworks.com/)**

### Resource Documentation

| Resource | Description |
|----------|-------------|
| `e2e_node` | Virtual machine instances |
| `e2e_vpc` | Virtual Private Cloud networks |
| `e2e_loadbalancer` | Load balancers for distributing traffic |
| `e2e_blockstorage` | Persistent block storage volumes |
| `e2e_objectstore` | S3-compatible object storage |
| `e2e_sfs` | Shared file storage (NFS) |
| `e2e_kubernetes` | Managed Kubernetes clusters |
| `e2e_dbaas_postgresql` | Managed PostgreSQL databases |
| `e2e_dbaas_mysql` | Managed MySQL databases |
| `e2e_dbaas_mariadb` | Managed MariaDB databases |
| `e2e_faas_function` | Serverless functions ⚡ |
| `e2e_scaler_group` | Auto-scaling groups |
| `e2e_ssh_key` | SSH key management |
| `e2e_security_groups` | Network security groups |
| `e2e_reserved_ip` | Reserved IP addresses |
| `e2e_container_registry` | Private container image registry |
| `e2e_image` | Custom VM images |

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.18+ is required).

### Testing

```bash
# Run unit tests
go test ./... -v

# Run acceptance tests (requires E2E credentials)
TF_ACC=1 go test ./... -v -timeout 120m
```

### Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

For more details, see [CONTRIBUTING.md](CONTRIBUTING.md) (if available).

## Examples

Check out the `/examples` directory (if available) for complete infrastructure examples including:

- Multi-tier web applications
- High-availability database setups
- Kubernetes cluster deployments
- Serverless function deployments
- Auto-scaling configurations

## Support

- **Issues**: Report bugs or request features via [GitHub Issues](https://github.com/e2eterraformprovider/terraform-provider-e2e/issues)
- **Documentation**: [E2E Networks Documentation](https://docs.e2enetworks.com/)
- **Support Portal**: [E2E Networks Support](https://cloud.e2enetworks.com/)

## License

This provider is distributed under the MIT License. See [LICENSE](LICENSE) for more information.

---

<div align="center">

**Built with ❤️ using E2E Networks Cloud**

[E2E Console](https://cloud.e2enetworks.com/) • [Documentation](https://docs.e2enetworks.com/) • [Terraform Provider](https://registry.terraform.io/providers/e2eterraformprovider/e2e/latest)

MIT License - See [LICENSE](LICENSE) file for details

</div>
