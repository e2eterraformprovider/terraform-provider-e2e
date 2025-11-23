<div align="center">
  <img src="https://e2enetworks.com/OnlyE2E.svg" alt="E2E Networks Logo" width="180"/>

# 🚀 Terraform Provider for E2E Networks

  <p>
    <a href="https://myaccount.e2enetworks.com/"><img src="https://img.shields.io/badge/E2E_Cloud-Console-blue?style=for-the-badge" alt="E2E Cloud Console"/></a>
    <a href="https://docs.e2enetworks.com/"><img src="https://img.shields.io/badge/E2E-Documentation-green?style=for-the-badge" alt="E2E Documentation"/></a>
    <a href="https://registry.terraform.io/providers/e2eterraformprovider/e2e/latest"><img src="https://img.shields.io/badge/Terraform-Provider-purple?style=for-the-badge" alt="Terraform E2E Provider"/></a>
    <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-yellow?style=for-the-badge" alt="License: MIT"/></a>
  </p>

**A Production-Ready Infrastructure as Code Solution for E2E Networks**

</div>

---

## Overview

The E2E Networks Terraform Provider enables you to manage your E2E Networks cloud infrastructure using Terraform's infrastructure-as-code approach. Build, modify, and version your cloud infrastructure safely and efficiently.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.13+
- E2E Networks account with API credentials

## Contributing to the Provider

### Prerequisites

- [Go](https://golang.org/doc/install) 1.24+
- [Terraform](https://www.terraform.io/downloads.html) 0.13+

### Clone and Build

```bash
# Set GOPATH (add to ~/.bashrc or ~/.zshrc)
export GOPATH=~/code/go

# Create directory structure and clone
mkdir -p $GOPATH/src/github.com/e2eterraformprovider
cd $GOPATH/src/github.com/e2eterraformprovider
git clone https://github.com/e2eterraformprovider/terraform-provider-e2e
cd terraform-provider-e2e

# Build the provider
make build
```

### Installing Locally

Install the provider for local development and testing:

```bash
# Set version and platform
VERSION="0.1.0"
PLATFORM="darwin_arm64"  # darwin_amd64 (Intel Mac) | darwin_arm64 (Apple Silicon) | linux_amd64 | windows_amd64

# Install to Terraform plugin directory
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/e2eterraformprovider/e2e/${VERSION}/${PLATFORM}
cp terraform-provider-e2e ~/.terraform.d/plugins/registry.terraform.io/e2eterraformprovider/e2e/${VERSION}/${PLATFORM}/
```

### Development Workflow

```bash
make fmt      # Format code
make lint     # Run linters
make test     # Run tests
make vendor   # Download and vendor dependencies (optional, creates vendor/ directory)
```

**Note**: The `vendor/` directory is gitignored. Dependencies are managed via `go.mod` and `go.sum`. Use `make vendor` only if you need local vendoring for offline/air-gapped builds or testing.

## Documentation

For detailed documentation on all resources and data sources, visit:

- **[Official Provider Documentation](https://registry.terraform.io/providers/e2eterraformprovider/e2e/latest/docs)**
- **[E2E Networks API Documentation](https://docs.e2enetworks.com/api/myaccount/)**
- **[E2E Networks Cloud Console](https://cloud.e2enetworks.com/)**

---

<div align="center">

**Built with ❤️ using E2E Networks Cloud**

[E2E Console](https://myaccount.e2enetworks.com/) • [Documentation](https://docs.e2enetworks.com/) • [Terraform Provider](https://registry.terraform.io/providers/e2eterraformprovider/e2e/latest)

MIT License - See [MIT License](https://opensource.org/licenses/MIT) for details

</div>
