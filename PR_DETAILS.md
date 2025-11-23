# Pull Request: Add comprehensive unit tests for SSH key, security group, and reserved IP

## Summary
This PR adds comprehensive TestAcc test cases for `e2e/ssh_key`, `e2e/security_group`, and `e2e/reserve_ip` following the DigitalOcean provider testing patterns as requested.

## Changes

### SSH Key Tests (`e2e/ssh_key/`)
- ✅ **resource_ssh_key_test.go** - Resource tests with 7 test cases:
  - Basic CRUD operations
  - Update restrictions (label and ssh_key cannot be updated)
  - Missing required arguments validation
  - Invalid SSH key format handling
  - Import functionality
- ✅ **datasource_ssh_key_test.go** - Data source tests with 4 test cases:
  - Single SSH key lookup
  - Non-existent SSH key handling
  - List all SSH keys
  - Missing required arguments validation
- ✅ **sweep.go** - Test cleanup for SSH keys

### Security Group Tests (`e2e/security_group/`)
- ✅ **resource_security_groups_test.go** - Resource tests with 9 test cases:
  - Basic security group creation with rules
  - Multiple rules handling (Inbound/Outbound)
  - Update operations (rules and default status)
  - Different network types (any, manual, myNetwork)
  - Invalid input validation (rule_type, protocol_name, network)
  - Missing required arguments validation
  - Import functionality
- ✅ **datasource_security_group_test.go** - Data source tests with 3 test cases:
  - Single security group lookup
  - Multiple rules validation
  - Non-existent security group handling
  - Missing required arguments validation
- ✅ **sweep.go** - Test cleanup for security groups

### Reserved IP Tests (`e2e/reserve_ip/`)
- ✅ **resource_reserve_ip_test.go** - Resource tests with 4 test cases:
  - Basic reserved IP creation and validation
  - Status checking
  - Missing required arguments validation
  - Import functionality
- ✅ **datasource_reserve_ips_test.go** - Data source tests with 3 test cases:
  - List all reserved IPs
  - List with created reserved IP
  - Missing required arguments validation
- ✅ **sweep.go** - Test cleanup for reserved IPs

## Test Coverage
All tests follow the DigitalOcean provider patterns including:
- ✅ Argument validation (required fields)
- ✅ Input validation (valid values for enums)
- ✅ Resource lifecycle (create, read, update, delete)
- ✅ Import/export functionality
- ✅ Data source functionality
- ✅ Error handling and edge cases
- ✅ Sweep functions for test resource cleanup

## Files Changed
- 9 new files created
- ~1,928 lines of test code added

## Testing Notes
- `make fmt` - ✅ Completed successfully
- `make vet` - ⚠️ Skipped due to network errors (offline check required)
- `make lint` - ⚠️ Skipped due to network errors (offline check required)
- `make test` - ⚠️ Skipped due to network errors (offline check required)

## Abstraction Evaluation
**Why we don't need helper files like droplets.go:**

Unlike DigitalOcean's droplet resource which has:
- Complex objects with many attributes
- Shared pagination logic
- Complex network address handling (IPv4/IPv6)
- Schema reuse across multiple resources

Our resources (SSH keys, security groups, reserved IPs) are:
- Simpler with fewer attributes
- Already have inline flattening where needed (e.g., `flattenSshKeys`, `flattenReserveIps`)
- Don't require pagination
- Less complex networking logic

The current structure is clean, follows existing patterns in the codebase, and doesn't require additional abstraction layers.

## Request
Please run `make vet`, `make lint`, and `make test` offline to verify code quality and test correctness.

---

## How to Create the PR

Visit this URL to create the pull request:
https://github.com/indykish/terraform-provider-e2e/pull/new/claude/add-e2e-unit-tests-01NCup6MiKus1QandsSipQqq

Or use the following command:
```bash
gh pr create --title "Add comprehensive unit tests for SSH key, security group, and reserved IP" --body-file PR_DETAILS.md
```
