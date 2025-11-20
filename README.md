<div align="center">
  <img src="https://e2enetworks.com/OnlyE2E.svg" alt="E2E Networks Logo" width="180"/>

  # 🚀 Terraform Provider for E2E Networks

  <p>
    <a href="https://cloud.e2enetworks.com/"><img src="https://img.shields.io/badge/E2E_Cloud-Console-blue?style=for-the-badge" alt="E2E Cloud Console"/></a>
    <a href="https://docs.e2enetworks.com/"><img src="https://img.shields.io/badge/E2E-Documentation-green?style=for-the-badge" alt="E2E Documentation"/></a>
    <a href="https://registry.terraform.io/providers/e2eterraformprovider/e2e/latest"><img src="https://img.shields.io/badge/Terraform-Provider-purple?style=for-the-badge" alt="Terraform E2E Provider"/></a>
    <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-yellow?style=for-the-badge" alt="License: MIT"/></a>
  </p>

  **A Production-Ready Infrastructure as Code Solution for E2E Networks Cloud**
</div>

---

## Overview

The E2E Networks Terraform Provider enables you to manage your E2E Networks cloud infrastructure using Terraform's infrastructure-as-code approach. Build, modify, and version your cloud infrastructure safely and efficiently.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.13+
- [Go](https://golang.org/doc/install) 1.21+ (for building the provider plugin)
- E2E Networks account with API credentials

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

To use the provider locally, copy it to Terraform's plugin directory:

```bash
# Set your desired version (use any version like 0.1.0, 1.0.0, etc.)
VERSION="0.1.0"
PLATFORM="linux_amd64"  # See platform options below

# Create Terraform's local plugin directory structure
# This tells Terraform where to find your locally built provider
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/e2eterraformprovider/e2e/${VERSION}/${PLATFORM}

# Copy the built provider binary to the plugin directory
cp terraform-provider-e2e ~/.terraform.d/plugins/registry.terraform.io/e2eterraformprovider/e2e/${VERSION}/${PLATFORM}/
```

**Platform options** - Set `PLATFORM` to match your operating system:
- Linux: `linux_amd64`
- macOS (Intel): `darwin_amd64`
- macOS (Apple Silicon): `darwin_arm64`
- Windows: `windows_amd64`

**Note**: The version number (e.g., `0.1.0`) can be any value you choose for local development.

## Documentation

For detailed documentation on all resources and data sources, visit:

- **[Official Provider Documentation](https://registry.terraform.io/providers/e2eterraformprovider/e2e/latest/docs)**
- **[E2E Networks API Documentation](https://docs.e2enetworks.com/api/myaccount/)**
- **[E2E Networks Cloud Console](https://cloud.e2enetworks.com/)**

## Support

- **Documentation**: [E2E Networks Documentation](https://docs.e2enetworks.com/)
- **Support Portal**: [E2E Networks Support](https://cloud.e2enetworks.com/)

---

<div align="center">

**Built with ❤️ using E2E Networks Cloud**

[E2E Console](https://cloud.e2enetworks.com/) • [Documentation](https://docs.e2enetworks.com/) • [Terraform Provider](https://registry.terraform.io/providers/e2eterraformprovider/e2e/latest)

MIT License - See [MIT License](https://opensource.org/licenses/MIT) for details

</div>
