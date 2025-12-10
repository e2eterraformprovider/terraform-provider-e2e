<div align="center">
  <img src="https://e2enetworks.com/OnlyE2E.svg" alt="E2E Networks Logo" width="180"/>

# 🚀 Terraform Provider for E2E Networks

<p>
  <a href="https://myaccount.e2enetworks.com/">
    <img src="https://img.shields.io/badge/E2E_Networks-Console-blue?style=for-the-badge" />
  </a>
  <a href="https://docs.e2enetworks.com/">
    <img src="https://img.shields.io/badge/E2E_Networks-Documentation-green?style=for-the-badge" />
  </a>
  <a href="https://registry.terraform.io/providers/e2eterraformprovider/e2e/latest">
    <img src="https://img.shields.io/badge/Terraform-Provider-purple?style=for-the-badge" />
  </a>
  <a href="https://opensource.org/licenses/MIT">
    <img src="https://img.shields.io/badge/License-MIT-yellow?style=for-the-badge" />
  </a>
</p>

**A production-ready Infrastructure-as-Code provider for managing E2E Networks Cloud resources.**

</div>

---

## Overview

The **E2E Networks Terraform Provider** enables you to provision, manage, and automate E2E Cloud infrastructure using Terraform’s Infrastructure-as-Code workflow. It supports programmatic creation and lifecycle management of nodes, volumes, networking, and other cloud resources.

---

## Requirements

- **Terraform 0.13+**
- **Go 1.24+** (required for development)
- **E2E Networks account** and valid **API credentials**

---

## Contributing

Contributions are welcome! Whether reporting bugs, improving documentation, or submitting new features, we appreciate your help.

Before submitting a pull request:

1. Follow the instructions in **Developing the Provider**
2. Ensure all tests pass (`make test` and `make testacc`)
3. Run the linters (`make lint`)
4. Provide a clear description of your changes

---

## Developing the Provider

This section describes how to set up your development environment, make changes, and build the provider locally.

### Prerequisites

Install the following tools:

- [Go](https://golang.org/doc/install) **1.24+**
- [Terraform](https://www.terraform.io/downloads.html) **0.13+**
- Git

Configure your Go workspace:

```bash
# Add to ~/.bashrc or ~/.zshrc
export GOPATH=~/code/go or <your_go_workspace_dir>
export PATH=$PATH:$GOPATH/bin
```

### Setup

Clone the repository:

```bash
mkdir -p $GOPATH/src/github.com/e2eterraformprovider
cd $GOPATH/src/github.com/e2eterraformprovider
git clone https://github.com/e2eterraformprovider/terraform-provider-e2e
cd terraform-provider-e2e
```

Build the provider (this will automatically download and tidy dependencies):

```bash
make build
```

### Development Workflow

1. Create a feature branch:

```bash
git checkout -b feature/your-feature-name
```

2. Make your changes.

3. Format and lint your code:

```bash
make fmt
make lint
```

4. Build and install the provider locally:

```bash
make install   # Installs into Terraform plugin directory
```

---

## Testing

The provider includes both unit tests and acceptance tests.

### Unit Tests

Run all unit tests (no cloud access required):

```bash
make test
```

### Acceptance Tests

These tests create real resources on E2E Cloud and **require valid API credentials**.

Set required environment variables:

```bash
export E2E_API_KEY="your-api-key"
export E2E_AUTH_TOKEN="your-auth-token"
export E2E_TEST_PROJECT_ID="your-project-id"
export E2E_TEST_REGION="your-region"
```

Run all acceptance tests:

```bash
make testacc
```

Run a specific test:

```bash
make testacc TESTARGS="-run TestAccE2ENode_Basic"
```

> ⚠️ **Warning:** Acceptance tests create real resources and may incur charges. Use a dedicated test project.

---

## Building the Provider

To build and install the provider locally for Terraform:

```bash
make install
```

This builds the binary in the repository root and installs it to Terraform's plugin directory. Terraform will pick it up automatically from there.

---

## Documentation

Official documentation for provider resources and data sources:

- **Terraform Provider Docs**  
  https://registry.terraform.io/providers/e2eterraformprovider/e2e/latest/docs

- **E2E Networks API Documentation**  
  https://docs.e2enetworks.com/api/myaccount/

- **E2E Cloud Console**  
  https://myaccount.e2enetworks.com/

---

<div align="center">

**Built with ❤️ on E2E Networks Cloud**

[E2E Console](https://myaccount.e2enetworks.com/) •  
[Documentation](https://docs.e2enetworks.com/) •  
[Terraform Provider](https://registry.terraform.io/providers/e2eterraformprovider/e2e/latest)

<br/>

MIT License — See the [MIT License](https://opensource.org/licenses/MIT) for full details.

</div>
