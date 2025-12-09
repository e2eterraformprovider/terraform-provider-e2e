---
page_title: "v3.0 Upgrade Guide"
subcategory: "Guides"
description: |-
  Complete guide for upgrading from v2.2.7 to v3.0.0 of the E2E Terraform Provider
---

# E2E Provider v3.0 Upgrade Guide

This guide provides comprehensive instructions for migrating from **v2.2.7 to v3.0.0**.

> **⚠️ Important**: If you prefer to continue using the existing schema, stay on **v2.2.7** (latest v2.x). Version pinning in Terraform allows you to control when you migrate.

---

## 📑 Table of Contents

### Quick Start
- [Overview](#overview)
- [Should I Upgrade?](#should-i-upgrade)
- [Quick Migration (TL;DR)](#quick-migration-tldr)

### Migration Guide
- [Breaking Changes Summary](#breaking-changes-summary)
- [Detailed Breaking Changes](#detailed-breaking-changes)
- [Step-by-Step Migration](#step-by-step-migration)
- [Automated Migration Tools](#automated-migration-tools)

### Support
- [Troubleshooting](#troubleshooting)
- [Getting Help](#getting-help)
- [Rollback Instructions](#rollback-instructions)

---

## Overview

### What's in v3.0.0?

Version 3.0.0 is a **major release** focused on schema consistency and standardization across all resources.

**Goals**:
- ✨ Consistent field naming across all resources
- 🎯 Alignment with industry standards (AWS, Azure, GCP patterns)
- 🧹 Elimination of field naming inconsistencies
- 📦 Single, clear way to specify each parameter

**What's NOT changing**:
- ✅ No changes to provider functionality
- ✅ No changes to API behavior
- ✅ No changes to resource capabilities
- ✅ Only field names are being standardized

### Migration Time

- **With automated script**: 15-30 minutes
- **Manual migration**: 1-2 hours (depending on codebase size)

---

## Should I Upgrade?

### ✅ Upgrade to v3.0.0 if:

- Starting a new project
- Want consistent, standardized schema
- Ready to adopt industry-standard naming conventions
- Have time to test migration thoroughly

### 🔒 Stay on v2.2.7 if:

- Running production workloads not ready for changes
- Don't have time to test migration now
- Want to delay migration to a maintenance window
- Prefer stability over new standardization

### 📌 Version Pinning Example

```hcl
# Stay on v2.2.7 (recommended for stable production)
terraform {
  required_providers {
    e2e = {
      source  = "e2enetworks/e2e"
      version = "2.2.7"  # Pin to specific version
    }
  }
}

# Or use v3.0.0 when ready
terraform {
  required_providers {
    e2e = {
      source  = "e2enetworks/e2e"
      version = "~> 3.0"  # Use v3.x
    }
  }
}
```

---

## Quick Migration (TL;DR)

For experienced users who want to migrate quickly:

### 1. Backup your configuration
```bash
git commit -am "Pre-v3 migration backup"
```

### 2. Run automated migration script
```bash
curl -o migrate-to-v3.sh https://raw.githubusercontent.com/e2enetworks/terraform-provider-e2e/main/scripts/migrate-to-v3.sh
chmod +x migrate-to-v3.sh
./migrate-to-v3.sh
```

### 3. Review and test changes
```bash
git diff  # Review changes
terraform fmt
terraform validate
terraform plan  # Should show no unexpected changes
```

### 4. Update provider version
```hcl
terraform {
  required_providers {
    e2e = {
      source  = "e2enetworks/e2e"
      version = "~> 3.0"
    }
  }
}
```

### 5. Apply
```bash
terraform init -upgrade
terraform plan
terraform apply
```

**Need details?** Continue reading for comprehensive migration instructions.

---

## Breaking Changes Summary

### 📊 All Changes at a Glance

| Change | Old | New | Resources | Impact |
|--------|-----|-----|-----------|--------|
| **Parameter** | `location` | `region` | All resources | 🔴 High |
| **Type** | `project_id` (int) | `project_id` (string) | Block Storage, K8s | 🟡 Low |
| **Parameter** | `plan_name` | `plan` | LB, ASG, MariaDB | 🟡 Medium |
| **Parameter** | `plan_id`/`sku_id` (input) | `plan` (computed) | Autoscaling | 🟡 Medium |
| **Attribute** | `creation_time` | `created_at` | Image | 🟢 Low |
| **Attribute** | `created_on` | `created_at` | Object Store | 🟢 Low |
| **Attribute** | `setup_status` | `status` | Container Registry | 🟢 Low |

### 🎯 Changes by Resource Type

**All Resources**: `location` → `region`

**Specific Resources**:
- **e2e_blockstorage**: `project_id` type change, `location` → `region`
- **e2e_kubernetes**: `project_id` type change, `location` → `region`
- **e2e_loadbalancer**: `plan_name` → `plan`, `location` → `region`
- **e2e_autoscaling**: `plan_name` → `plan`, remove `plan_id`/`sku_id` inputs, `location` → `region`
- **e2e_dbaas_mariadb**: `plan_name` → `plan`, `location` → `region`
- **e2e_objectstore**: `created_on` → `created_at`, `location` → `region`
- **e2e_image**: `creation_time` → `created_at`, `location` → `region`
- **e2e_container_registry**: `setup_status` → `status`, `location` → `region`

---

## Detailed Breaking Changes

### 1. Region/Location Parameter

**Impact**: 🔴 **High** - Affects ALL resources

**Change**:
```hcl
# ❌ v2.2.7
resource "e2e_node" "example" {
  name     = "my-node"
  location = "Delhi NCR"
}

# ✅ v3.0.0
resource "e2e_node" "example" {
  name   = "my-node"
  region = "Delhi NCR"
}
```

**Why?**
- Industry standard (AWS, Azure, GCP all use `region`)
- Reduces confusion between "location" and "region"
- More intuitive for users familiar with other cloud providers

**Automated fix**: ✅ Yes

---

### 2. Project ID Type

**Impact**: 🟡 **Low** - Auto-converts in most cases

**Change**:
```hcl
# ❌ v2.2.7 (works but not recommended)
resource "e2e_blockstorage" "example" {
  project_id = 12345
}

# ✅ v3.0.0 (recommended)
resource "e2e_blockstorage" "example" {
  project_id = "12345"
}
```

**Why?**
- Strings provide better validation
- Prevents type confusion
- Allows for future extensibility

**Note**: Terraform auto-converts integers to strings, so existing configs often work without changes.

**Automated fix**: ✅ Yes

---

### 3. Plan Field Standardization

**Impact**: 🟡 **Medium** - Three resources affected

**Change**:
```hcl
# ❌ v2.2.7
resource "e2e_loadbalancer" "lb" {
  plan_name = "E2E-LB-2"
}

resource "e2e_autoscaling" "asg" {
  plan_name = "c4.large.x86"
}

resource "e2e_dbaas_mariadb" "db" {
  plan_name = "mariadb-standard-1"
}

# ✅ v3.0.0
resource "e2e_loadbalancer" "lb" {
  plan = "E2E-LB-2"
}

resource "e2e_autoscaling" "asg" {
  plan = "c4.large.x86"
}

resource "e2e_dbaas_mariadb" "db" {
  plan = "mariadb-standard-1"
}
```

**Why?**
- Consistency with `e2e_node`, `e2e_dbaas_mysql`, `e2e_dbaas_postgresql`
- Shorter, clearer naming
- Follows Terraform provider conventions

**Automated fix**: ✅ Yes

---

### 4. Plan/SKU Consolidation (Autoscaling)

**Impact**: 🟡 **Medium** - Autoscaling only

**Change**:
```hcl
# ❌ v2.2.7 (multiple ways - confusing!)
resource "e2e_autoscaling" "example" {
  plan_id = "12345"    # Option 1
  # OR
  sku_id = "12345"     # Option 2
  # OR
  plan_name = "c4.large.x86"  # Option 3
}

# ✅ v3.0.0 (one clear way)
resource "e2e_autoscaling" "example" {
  plan = "c4.large.x86"  # Only way
}

# Note: plan_id and sku_id still available as outputs
output "details" {
  value = {
    plan_id = e2e_autoscaling.example.plan_id  # ✅ Works
    sku_id  = e2e_autoscaling.example.sku_id   # ✅ Works
  }
}
```

**Why?**
- Eliminates confusion from multiple input methods
- Clearer user experience
- IDs are implementation details, users specify by name

**Automated fix**: ⚠️ Partial (removes `plan_name`, but manual review needed for `plan_id`/`sku_id` removal)

---

### 5. Timestamp Field (Image)

**Impact**: 🟢 **Low** - Output only

**Change**:
```hcl
# ❌ v2.2.7
output "when_created" {
  value = e2e_image.example.creation_time
}

# ✅ v3.0.0
output "when_created" {
  value = e2e_image.example.created_at
}
```

**Why?**
- Standard naming across all resources
- Matches API field naming

**Automated fix**: ✅ Yes

---

### 6. Timestamp Field (Object Store)

**Impact**: 🟢 **Low** - Output only

**Change**:
```hcl
# ❌ v2.2.7
output "bucket_created" {
  value = e2e_objectstore.bucket.created_on
}

data "e2e_objectstore" "existing" {
  name = "my-bucket"
}

output "data_created" {
  value = data.e2e_objectstore.existing.created_on
}

# ✅ v3.0.0
output "bucket_created" {
  value = e2e_objectstore.bucket.created_at
}

data "e2e_objectstore" "existing" {
  name = "my-bucket"
}

output "data_created" {
  value = data.e2e_objectstore.existing.created_at
}
```

**Why?**
- Consistency with all other resources
- Aligns with API field naming

**Automated fix**: ✅ Yes

---

### 7. Status Field (Container Registry)

**Impact**: 🟢 **Low** - Output only

**Change**:
```hcl
# ❌ v2.2.7
output "registry_status" {
  value = e2e_container_registry.reg.setup_status
}

# ✅ v3.0.0
output "registry_status" {
  value = e2e_container_registry.reg.status
}
```

**Why?**
- Simpler, clearer naming
- "status" is more standard than "setup_status"

**Automated fix**: ✅ Yes

---

## Step-by-Step Migration

### Step 1: Backup Current Configuration

```bash
# Commit current state
git add .
git commit -m "Pre-v3.0 migration backup"

# Or create a backup
cp -r . ../terraform-backup
```

### Step 2: Choose Migration Method

**Option A: Automated Script (Recommended)**
- ✅ Fast (5-10 minutes)
- ✅ Handles most changes automatically
- ✅ Creates backups
- ⚠️ Requires manual review

[Jump to Automated Migration →](#automated-migration-tools)

**Option B: Manual Migration**
- ✅ Full control
- ✅ Learn each change
- ⏱️ Takes longer
- ✅ Good for small codebases

[Jump to Manual Migration →](#manual-migration)

### Step 3: Test Changes

```bash
# Format code
terraform fmt -recursive

# Validate syntax
terraform validate

# Review plan
terraform plan
```

**What to check**:
- No syntax errors
- No unexpected resource replacements
- No unexpected changes to infrastructure
- Plan shows only attribute name updates (if any)

### Step 4: Update Provider Version

```hcl
terraform {
  required_providers {
    e2e = {
      source  = "e2enetworks/e2e"
      version = "~> 3.0"  # Update to v3
    }
  }
}
```

### Step 5: Apply Changes

```bash
# Download v3 provider
terraform init -upgrade

# Review final plan
terraform plan

# Apply (state updates only, no infrastructure changes)
terraform apply

# Refresh state to update attribute names
terraform refresh
```

---

## Automated Migration Tools

### Migration Script

This script automatically updates your `.tf` files with all v3.0.0 changes:

**Features**:
- ✅ Automatically creates backups (`.pre-v3-backup` extension)
- ✅ Handles all 7 breaking changes
- ✅ Progress indicators
- ✅ Rollback support
- ⚠️ Requires manual review afterward

### Download and Run

```bash
# Download script
curl -o migrate-to-v3.sh https://raw.githubusercontent.com/e2enetworks/terraform-provider-e2e/main/scripts/migrate-to-v3.sh

# Make executable
chmod +x migrate-to-v3.sh

# Run (creates backups automatically)
./migrate-to-v3.sh
```

### Script Content

Save as `migrate-to-v3.sh`:

```bash
#!/bin/bash
# E2E Provider v2.2.7 → v3.0.0 Migration Script

set -e

echo "🚀 E2E Provider v3.0.0 Migration"
echo "=================================="
echo ""
echo "This script will update your .tf files for v3.0.0 compatibility."
echo ""

# Confirm
read -p "Continue? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ Cancelled"
    exit 1
fi

# Backup
echo ""
echo "📦 Step 1/7: Creating backups..."
find . -name "*.tf" -type f -exec cp {} {}.pre-v3-backup \;
echo "✅ Backups created (.pre-v3-backup extension)"

# Migration 1: location → region
echo ""
echo "🔄 Step 2/7: Migrating location → region..."
find . -name "*.tf" -type f -exec sed -i.bak 's/\blocation\s*=/region =/g' {} \;
echo "✅ Complete"

# Migration 2: Quote project_id
echo ""
echo "🔄 Step 3/7: Quoting project_id values..."
find . -name "*.tf" -type f -exec sed -i.bak 's/project_id\s*=\s*\([0-9]\+\)/project_id = "\1"/g' {} \;
echo "✅ Complete"

# Migration 3: plan_name → plan
echo ""
echo "🔄 Step 4/7: Migrating plan_name → plan..."
find . -name "*.tf" -type f -exec sed -i.bak 's/\bplan_name\s*=/plan =/g' {} \;
echo "✅ Complete"

# Migration 4: creation_time → created_at
echo ""
echo "🔄 Step 5/7: Migrating creation_time → created_at..."
find . -name "*.tf" -type f -exec sed -i.bak 's/\.creation_time/.created_at/g' {} \;
echo "✅ Complete"

# Migration 5: created_on → created_at
echo ""
echo "🔄 Step 6/7: Migrating created_on → created_at..."
find . -name "*.tf" -type f -exec sed -i.bak 's/\.created_on/.created_at/g' {} \;
echo "✅ Complete"

# Migration 6: setup_status → status
echo ""
echo "🔄 Step 7/7: Migrating setup_status → status..."
find . -name "*.tf" -type f -exec sed -i.bak 's/\.setup_status/.status/g' {} \;
echo "✅ Complete"

# Cleanup
echo ""
echo "🧹 Cleaning up temporary files..."
find . -name "*.tf.bak" -delete
echo "✅ Complete"

echo ""
echo "✅ Migration Complete!"
echo ""
echo "⚠️  IMPORTANT: Manual Review Required"
echo "======================================"
echo ""
echo "1. Review changes:"
echo "   git diff"
echo ""
echo "2. Check for e2e_autoscaling resources:"
echo "   - Remove 'plan_id' and 'sku_id' input parameters"
echo "   - Keep them only as outputs"
echo ""
echo "3. Test your configuration:"
echo "   terraform fmt"
echo "   terraform validate"
echo "   terraform plan"
echo ""
echo "4. Update provider version in your configuration:"
echo "   version = \"~> 3.0\""
echo ""
echo "5. Initialize and apply:"
echo "   terraform init -upgrade"
echo "   terraform plan"
echo "   terraform apply"
echo ""
echo "💾 Backups saved with .pre-v3-backup extension"
echo ""
echo "❓ Need help? See docs/guides/v3-upgrade-guide.md"
```

### Post-Script Manual Steps

After running the script, you **must** manually:

1. **Review all changes**:
   ```bash
   git diff
   ```

2. **Check Autoscaling resources** for `plan_id`/`sku_id` usage:
   ```bash
   grep -r "plan_id\s*=" *.tf
   grep -r "sku_id\s*=" *.tf
   ```
   - Remove these as **input parameters**
   - Keep as **outputs** only

3. **Test thoroughly**:
   ```bash
   terraform fmt
   terraform validate
   terraform plan  # Review carefully!
   ```

---

## Manual Migration

If you prefer to migrate manually:

### 1. Find All Occurrences

```bash
# Find location usage
grep -rn "location\s*=" *.tf

# Find plan_name usage
grep -rn "plan_name\s*=" *.tf

# Find timestamp attributes
grep -rn "\.creation_time" *.tf
grep -rn "\.created_on" *.tf
grep -rn "\.setup_status" *.tf

# Find autoscaling plan_id/sku_id
grep -rn "plan_id\s*=" *.tf
grep -rn "sku_id\s*=" *.tf
```

### 2. Replace in Your Editor

Use your editor's find & replace:

| Find | Replace |
|------|---------|
| `location =` | `region =` |
| `plan_name =` | `plan =` |
| `.creation_time` | `.created_at` |
| `.created_on` | `.created_at` |
| `.setup_status` | `.status` |

### 3. Manual Fixes

For Autoscaling resources, manually:
- Remove `plan_id = "..."` lines (as input)
- Remove `sku_id = "..."` lines (as input)
- Add `plan = "plan-name"` (convert from ID to name)

---

## Troubleshooting

### Error: "Unknown parameter: location"

**Solution**: You missed replacing a `location` parameter.

```bash
# Find all remaining
grep -rn "location\s*=" *.tf

# Replace with region
```

### Error: "Unknown parameter: plan_name"

**Solution**: You missed replacing a `plan_name` parameter.

```bash
# Find all remaining
grep -rn "plan_name\s*=" *.tf

# Replace with plan
```

### Error: "Unknown attribute: created_on"

**Solution**: You have an output referencing old attribute.

```bash
# Find all remaining
grep -rn "\.created_on" *.tf

# Replace with .created_at
```

### Plan Shows Unexpected Resource Changes

**Solution**: Run `terraform refresh` to update state.

```bash
terraform refresh
terraform plan  # Should now show no changes
```

### Autoscaling Errors About plan_id/sku_id

**Solution**: Ensure you removed these as **inputs**, not outputs.

```hcl
# ❌ Wrong - as input
resource "e2e_autoscaling" "example" {
  plan_id = "12345"  # Remove this
}

# ✅ Correct - as output
output "plan_info" {
  value = e2e_autoscaling.example.plan_id  # Keep this
}
```

---

## Getting Help

### 1. Documentation
- [CHANGELOG.md](../../CHANGELOG.md) - All changes
- [Resource Docs](../) - Updated documentation

### 2. Check Existing Issues
Search: [v3-migration label](https://github.com/e2enetworks/terraform-provider-e2e/labels/v3-migration)

### 3. Open New Issue
1. Go to: https://github.com/e2enetworks/terraform-provider-e2e/issues
2. Click "New Issue"
3. Title: "v3 Migration: [brief description]"
4. Label: `v3-migration`
5. Include:
   - Current provider version
   - Error message
   - Configuration snippet
   - What you've tried

### 4. Community Support
- **GitHub Discussions**: [Start a discussion](https://github.com/e2enetworks/terraform-provider-e2e/discussions)

---

## Rollback Instructions

### Option 1: Stay on v2.2.7 (Recommended)

If you haven't upgraded yet or want to rollback:

```hcl
terraform {
  required_providers {
    e2e = {
      source  = "e2enetworks/e2e"
      version = "2.2.7"  # Pin to v2.2.7
    }
  }
}
```

```bash
terraform init -upgrade
terraform plan
```

### Option 2: Restore from Backups

If you used the migration script:

```bash
# Restore .tf files from backup
find . -name "*.tf.pre-v3-backup" -exec bash -c 'mv "$0" "${0%.pre-v3-backup}"' {} \;

# Downgrade provider
terraform init -upgrade

# Verify
terraform plan
```

### Option 3: Git Rollback

If you committed to git:

```bash
# View commits
git log --oneline -5

# Revert migration commit
git revert <commit-hash>

# Or reset (⚠️ destructive)
git reset --hard HEAD~1
```

---

## Version Support

### v3.0.0 (Current)
- ✅ Active development
- ✅ New features
- ✅ Bug fixes
- ✅ Security updates

### v2.2.7 (Stable)
- ✅ Critical bug fixes
- ✅ Security updates
- ⚠️ No new features
- 📅 Supported until: TBD

### Older v2.x (Legacy)
- ❌ No longer supported
- ⚠️ Upgrade to v2.2.7 or v3.0.0

---

## FAQ

### Q: Can I use both v2.2.7 and v3.0.0 in different workspaces?

**A**: Yes! Version pinning is per-configuration. Pin each workspace as needed:

```hcl
# Workspace 1 (production) - stays on v2.2.7
version = "2.2.7"

# Workspace 2 (staging) - tests v3.0.0
version = "~> 3.0"
```

### Q: Will v2.2.7 receive updates?

**A**: Yes, critical bug fixes and security updates. No new features.

### Q: Can I rollback after upgrading to v3.0.0?

**A**: Yes, see [Rollback Instructions](#rollback-instructions). Your infrastructure won't change, only state file needs updating.

### Q: Do I need to migrate all resources at once?

**A**: Yes, you cannot mix v2 and v3 schema in the same workspace. Use separate workspaces if needed.

### Q: Will this affect my actual infrastructure?

**A**: No! v3.0.0 only changes **schema names**, not API behavior. Your infrastructure remains unchanged.

---

## Additional Resources

- **[CHANGELOG.md](../../CHANGELOG.md)** - Detailed change history
- **[GitHub Issues](https://github.com/e2enetworks/terraform-provider-e2e/issues)** - Report issues
- **[GitHub Discussions](https://github.com/e2enetworks/terraform-provider-e2e/discussions)** - Ask questions

---

*Last updated: 2025-12-01*

**Questions?** [Open an issue](https://github.com/e2enetworks/terraform-provider-e2e/issues/new) with the `v3-migration` label.
